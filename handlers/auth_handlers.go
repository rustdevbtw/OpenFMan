package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Auth middleware for password protection
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check session or authentication
		session, err := c.Cookie("session_id")
		if err != nil || session != "valid" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Normalize the parameters by trimming the leading "/"
		for i, param := range c.Params {
			if len(param.Value) > 0 && param.Value[0] == '/' {
				c.Params[i].Value = param.Value[1:]
			}
		}

		c.Next()
	}
}

// Login handles simple password authentication
func Login(c *gin.Context) {
	var request struct {
		Password string `json:"password"`
	}
	if err := c.BindJSON(&request); err != nil {
		jsonError(c, http.StatusBadRequest, "Invalid request")
		return
	}

	if request.Password != viper.GetString("auth.password") {
		jsonError(c, http.StatusUnauthorized, "Invalid password")
		return
	}

	// Set cookie for authentication
	c.SetCookie("session_id", "valid", 3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// Logout handles user logout by clearing the session cookie
func Logout(c *gin.Context) {
	c.SetCookie("session_id", "", -1, "/", "", false, true) // Set max age to -1 to delete the cookie
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
