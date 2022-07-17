package types

// ErrorResponse Represents an error.
type ErrorResponse struct {

	// The error message.
	// Required: true
	Message string `json:"message"`

	// The error message.
	// Required: false
	StatusCode int `json:"statusCode"`
}
