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

// Error404 return true if the error is a RequestError and the status code is 404
func Error404(err error) bool {
	if err == nil {
		return false
	}
	serr := err.(*RequestError)
	return serr.Code == 404
}
