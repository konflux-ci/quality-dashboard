package types

// SuccessResponse Represents an success api request.
type SuccessResponse struct {

	// The message.
	// Required: true
	Message string `json:"message"`

	// The error message.
	// Required: false
	StatusCode int `json:"statusCode"`
}
