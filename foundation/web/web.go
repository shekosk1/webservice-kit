// Package web contains a small web framework extension to decouple the code from 3rd party libs.
package web

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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
		v := Values{
			TraceID: uuid.NewString(),
			Now:     time.Now().UTC(),
		}
		ctx := context.WithValue(r.Context(), key, &v)

		if err := handler(ctx, w, r); err != nil {
			return //WIP
		}

		//WIP
	}

	a.Mux.MethodFunc(method, path, h)
}
