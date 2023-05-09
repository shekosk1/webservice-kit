// Package metrics constructs the metrics that the app will track.
package metrics

import (
	"context"
	"expvar"
	"runtime"
)

// single instance of metrics.
var m *metrics

// metrics holds the set of metrics to be gathered.
// pkg expvar allows a concurrency safe work.
type metrics struct {
	goroutines *expvar.Int
	requests   *expvar.Int
	errors     *expvar.Int
	panics     *expvar.Int
}

// init constructs the metrics value that will be used to track and capture metrics.
// Thus, we can ensure the initialization only happens once.
func init() {
	m = &metrics{
		goroutines: expvar.NewInt("goroutines"),
		requests:   expvar.NewInt("requests"),
		errors:     expvar.NewInt("errors"),
		panics:     expvar.NewInt("panics"),
	}
}

type ctxKey int

const key ctxKey = 1

// Set sets the metrics data into the context.
func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, m)
}

// AddGoroutines refreshes the goroutine metric every 100 requests.
func AddGoroutine(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		if v.requests.Value()%100 == 0 {
			v.goroutines.Set(int64(runtime.NumGoroutine()))
		}
	}
}

// AddRequests increments the request metric by one.
func AddRequests(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.requests.Add(1)
	}
}

// AddErrors increments the errors metric by one.
func AddErrors(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.errors.Add(1)
	}
}

// AddPanics increments the panics metric by one.
func AddPanics(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.panics.Add(1)
	}
}
