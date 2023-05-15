// This program performs administrative tasks for the fairsplit service.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/ardanlabs/conf/v3"
	"github.com/shekosk1/webservice-kit/app/tooling/fairsplit-admin/commands"
	"github.com/shekosk1/webservice-kit/foundation/vault"

	"go.uber.org/zap"
)

var build = "develop"

// ErrHelp provides context that help was given.
var ErrHelp = errors.New("provided help")

type config struct {
	conf.Version
	Args  conf.Args
	Vault struct {
		KeysFolder string `conf:"default:/src/zarf/keys/"`
		Address    string `conf:"default:http://vault-service.fairsplit-system.svc.cluster.local:8200"`
		Token      string `conf:"default:mytoken,mask"`
		MountPath  string `conf:"default:secret"`
	}
}

func main() {
	if err := run(zap.NewNop().Sugar()); err != nil {
		if !errors.Is(err, ErrHelp) {
			fmt.Println("ERROR", err)
		}
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	cfg := config{
		Version: conf.Version{
			Build: build,
			Desc:  "web service admin tool kit",
		},
	}

	const prefix = "FAIRSPLIT"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}

		out, err := conf.String(&cfg)
		if err != nil {
			return fmt.Errorf("generating config for output: %w", err)
		}
		log.Infow("startup", "config", out)

		return fmt.Errorf("parsing config: %w", err)
	}

	return processCommands(cfg.Args, log, cfg)
}

// processCommands handles the execution of the commands specified on
// the command line.
func processCommands(args conf.Args, log *zap.SugaredLogger, cfg config) error {
	vaultConfig := vault.Config{
		Address:   cfg.Vault.Address,
		Token:     cfg.Vault.Token,
		MountPath: cfg.Vault.MountPath,
	}

	switch args.Num(0) {
	case "genkey":
		if err := commands.GenKey(); err != nil {
			return fmt.Errorf("key generation: %w", err)
		}

	case "vault":
		if err := commands.Vault(vaultConfig, cfg.Vault.KeysFolder); err != nil {
			return fmt.Errorf("setting private key: %w", err)
		}

	case "vault-init":
		if err := commands.VaultInit(vaultConfig); err != nil {
			return fmt.Errorf("initializing vault instance: %w", err)
		}

	default:
		fmt.Println("genkey:     generate a set of private/public key files")
		fmt.Println("vault:      load private keys into vault system")
		fmt.Println("vault-init: initialize a new vault instance")
		fmt.Println("provide a command to get more help.")
		return ErrHelp
	}

	return nil
}
