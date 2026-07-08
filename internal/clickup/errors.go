package clickup

import (
	"encoding/json"
	"errors"
	"fmt"
)

// APIError represents an error response from the ClickUp API, decoded from
// its {"err": "...", "ECODE": "..."} envelope.
type APIError struct {
	StatusCode int
	Err        string `json:"err"`
	ECode      string `json:"ECODE"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("clickup api error %d [%s]: %s", e.StatusCode, e.ECode, e.Err)
}

func decodeAPIError(statusCode int, body []byte) error {
	apiErr := &APIError{StatusCode: statusCode}
	if len(body) > 0 {
		_ = json.Unmarshal(body, apiErr)
	}
	if apiErr.Err == "" {
		apiErr.Err = string(body)
	}
	return apiErr
}

// IsNotFound reports whether err is a ClickUp APIError with a 404 status.
func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 404
	}
	return false
}

// IsRateLimited reports whether err is a ClickUp APIError with a 429 status.
func IsRateLimited(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 429
	}
	return false
}
