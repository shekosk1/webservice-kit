// Package v1 holds types used by the web app for v1.
package v1

import "errors"

// ErrorResponse is the structure used by the API to respond to the client
// when a failure happens.
type ErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

// RequestError is used to pass an error during the request through the
// application with the web specific context.
type RequestError struct {
	Err    error
	Status int
}

// Error implements the error interface.
func (re *RequestError) Error() string {
	return re.Err.Error()
}

// NewRequestError wraps the provided error and its http status, returning a RequestError.
func NewRequestError(err error, status int) error {
	return &RequestError{err, status}
}

// IsRequestError is a helper that checks if an error of type RequestError exists.
func IsRequestError(err error) bool {
	var re *RequestError
	return errors.As(err, &re)
}

// GetRequestError returns a copy of the RequestError pointer.
func GetRequestError(err error) *RequestError {
	var re *RequestError
	if !errors.As(err, &re) {
		return nil
	}

	return re
}
