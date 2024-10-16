package handlers

import (
	"github.com/gin-gonic/gin"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// jsonError sends a JSON error response
func jsonError(c *gin.Context, status int, message string) {
	c.JSON(status, ErrorResponse{Status: status, Message: message})
}
