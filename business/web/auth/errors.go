package auth

import (
	"errors"
	"fmt"
)

// authError is used to pass an error during the request through the
// application with auth specific context.
type authError struct {
	msg string
}

// Error implements the error interface.
func (ae *authError) Error() string {
	return ae.msg
}

// NewAuthError creates an AuthError for the provided message.
func NewAuthError(formatter string, args ...any) error {
	return &authError{
		msg: fmt.Sprintf(formatter, args...),
	}
}

// IsAuthError checks if an error of type AuthError exists.
func IsAuthError(err error) bool {
	var ae *authError
	return errors.As(err, &ae)
}
