package harbor

import (
	"net/http"

	"github.com/hashicorp/errwrap"
)

type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}

func ErrorIs404(err error) bool {
	harborError, ok := errwrap.GetType(err, &APIError{}).(*APIError)

	return ok && harborError != nil && harborError.Code == http.StatusNotFound
}
