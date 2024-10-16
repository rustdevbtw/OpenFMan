package models

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
