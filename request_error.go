package tvdb

import "fmt"

// RequestError is raised when a response from TVDB api return an http error code
// different from 200.
type RequestError struct {
	Code int
}

// Implement error interface method.
func (e *RequestError) Error() string {
	return fmt.Sprintf("Get a response with status code %d", e.Code)
}

// HaveCodeError return true if the param err is a RequestError and the status
// code is equal to the param code.
func HaveCodeError(code int, err error) bool {
	if err == nil {
		return false
	}
	serr := err.(*RequestError)
	return serr.Code == code
}
