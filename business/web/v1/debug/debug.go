// Package debug provides endpoints to debug the app.
package debug

import (
	"expvar"
	"net/http"
	"net/http/pprof"
)

// StandardLibraryMux bypass the use of the DefaultServerMux for security reasons.
// For example, a dependency could inject a handler into our service without us knowing it.
func StandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}
