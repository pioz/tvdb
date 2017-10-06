package tvdb

import "fmt"

// RequestError is raised when a response from TVDB api return an error code different from 200
type RequestError struct {
	Code int
}

// Error implementation
func (e *RequestError) Error() string {
	return fmt.Sprintf("Get a response with status code %d", e.Code)
}

// HaveCodeError return true if the error is a RequestError and the status code is 404
func HaveCodeError(code int, err error) bool {
	if err == nil {
		return false
	}
	serr := err.(*RequestError)
	return serr.Code == code
}
