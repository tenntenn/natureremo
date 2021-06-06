package natureremo

import (
	"fmt"
)

// APIError is an error from Nature Remo API.
type APIError struct {
	HTTPStatus int    `json:"-"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
}

var _ error = (*APIError)(nil)

func (err *APIError) Error() string {
	if err.Message == "" {
		return fmt.Sprintf("request failed with status code %d", err.HTTPStatus)
	}
	return fmt.Sprintf("StatusCode: %d, Code: %d, Message: %s", err.HTTPStatus, err.Code, err.Message)
}
