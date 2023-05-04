// Package handlers manages the Fairsplit API
package handlers

import (
	"net/http"
	"os"

	"github.com/shekosk1/webservice-kit/app/services/fairsplit-api/handlers/v1/testgrp"
	"github.com/shekosk1/webservice-kit/business/web/v1/mid"
	"github.com/shekosk1/webservice-kit/foundation/web"
	"go.uber.org/zap"
)

// APIMuxConfig contains the mandatory elements required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

// APIMux constructs an http handler with all the application routes.
func APIMux(cfg APIMuxConfig) *web.App {
	app := web.NewApp(cfg.Shutdown)

	app.Handle(http.MethodGet, "/status", testgrp.Status, mid.Logger(cfg.Log))
	app.Handle(http.MethodGet, "/empty", testgrp.Empty, mid.Logger(cfg.Log))

	return app

}
