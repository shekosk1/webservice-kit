package web

// Middleware is a type function designed to run before and/or after another Handler.
type Middleware func(Handler) Handler

// wrapMiddleware creates a new handler by wrapping a set of middlewares.
func wrapMiddleware(mw []Middleware, handler Handler) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}
