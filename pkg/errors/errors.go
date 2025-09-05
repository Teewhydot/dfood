package errors

import "fmt"

type HTTPError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func NewHTTPError(statusCode int, message string, err error) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}

func GetStatusCode(err error) (int, bool) {
	if httpErr, ok := err.(*HTTPError); ok {
		return httpErr.StatusCode, true
	}
	return 0, false
}

func GetErrorMessage(err error) (string, bool) {
	if httpErr, ok := err.(*HTTPError); ok {
		return httpErr.Message, true
	}
	return "", false
}