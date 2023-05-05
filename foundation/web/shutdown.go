package web

import (
	"errors"
)

// shutdownError helps to gracefully stop the service.
type shutdownError struct {
	Message string
}

// Error is the shutdownError implementation of the error interface.
func (se *shutdownError) Error() string {
	return se.Message
}

// NewShutdownError returns an error that causes the framework to signal a shutdown.
func NewShutdownError(message string) error {
	return &shutdownError{message}
}

// IsShutdown is a helper that checks if err contains a shutdownError.
func IsShutdown(err error) bool {
	var se *shutdownError

	return errors.As(err, &se)
}
