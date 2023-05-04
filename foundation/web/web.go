// Package web contains a small web framework extension to decouple the code from 3rd party libs.
package web

import (
	"context"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

// Handler handles an http request.
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint of the app. It configures the context object for each http handler.
type App struct {
	*chi.Mux
	shutdown chan os.Signal
	mw       []Middleware
}

// NewApp returns an App value that handles a set of routes for the app.
func NewApp(shudown chan os.Signal, mw ...Middleware) *App {
	return &App{
		Mux:      chi.NewMux(),
		shutdown: shudown,
		mw:       mw,
	}
}

// Handle associates a handler function with the specified http method and path.
func (a *App) Handle(method string, path string, handler Handler, mw ...Middleware) {
	handler = wrapMiddleware(mw, handler)
	handler = wrapMiddleware(a.mw, handler)

	h := func(w http.ResponseWriter, r *http.Request) {

		//WIP LOG
		if err := handler(r.Context(), w, r); err != nil {
			return //WIP
		}
		//WIP LOG
	}

	a.Mux.MethodFunc(method, path, h)
}
