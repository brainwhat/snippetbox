// We create out own errors and not relying on driver-specific errors
// For encapsulation
package models

import "errors"

var (
	ErrNoRecord           = errors.New("models: no matching recotd found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)
