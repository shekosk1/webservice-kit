package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/shekosk1/webservice-kit/app/services/fairsplit-api/handlers"
	"github.com/shekosk1/webservice-kit/business/web/auth"
	"github.com/shekosk1/webservice-kit/business/web/v1/debug"
	"github.com/shekosk1/webservice-kit/foundation/logger"
	"github.com/shekosk1/webservice-kit/foundation/vault"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"

	"github.com/ardanlabs/conf/v3"
)

var build = "develop"

/* add hashicorp vault to config and foundation/vault */

func main() {
	log, err := logger.New("FAIRSPLIT-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	/*==========================================================================
		SET GOMAXPROCS.
	==========================================================================*/

	opt := maxprocs.Logger(log.Infof)
	if _, err := maxprocs.Set(opt); err != nil {
		return fmt.Errorf("maxprocs: %w", err)
	}
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	/*==========================================================================
		Set app config.
	==========================================================================*/

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:180s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
		}
		Auth struct {
			Issuer string `conf:"default:fairsplit project"`
		}
		Vault struct {
			Address   string `conf:"default:http://vault-service.fairsplit-system.svc.cluster.local:8200"`
			MountPath string `conf:"default:secret"`
			Token     string `conf:"default:mytoken,mask"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "web service kit",
		},
	}

	const prefix = "FAIRSPLIT"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	/*==========================================================================
		Print config.
	==========================================================================*/

	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)
	expvar.NewString("build").Set(build)

	/*==========================================================================
		Start auth support.
	==========================================================================*/

	log.Infow("startup", "status", "initializing authentication support")

	vault, err := vault.New(vault.Config{
		Address:   cfg.Vault.Address,
		Token:     cfg.Vault.Token,
		MountPath: cfg.Vault.MountPath,
	})
	if err != nil {
		return fmt.Errorf("constructing vault: %w", err)
	}

	authCfg := auth.Config{
		Log:       log,
		KeyLookup: vault,
	}

	auth, err := auth.New(authCfg)
	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}
	/*==========================================================================
		Start debug service.
	==========================================================================*/

	log.Infow("startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.StandardLibraryMux()); err != nil {
			log.Errorw("shutdown", "status", "debug router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	/*==========================================================================
		Start API service.
	==========================================================================*/

	log.Infow("startup", "status", "initializing API service", "host", cfg.Web.APIHost)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	apiMux := handlers.APIMux(
		handlers.APIMuxConfig{
			Shutdown: shutdown,
			Log:      log,
			Auth:     auth,
		},
	)

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	listenServerErrors := make(chan error, 1)

	go func() {
		log.Infow("startup", "status", "API service started", "host", api.Addr)
		listenServerErrors <- api.ListenAndServe()
	}()

	/*==========================================================================
		App shutdown.
	==========================================================================*/

	select {
	case err := <-listenServerErrors:
		return fmt.Errorf("server error: %w", err)

	case stopSignal := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", stopSignal)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", stopSignal)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
