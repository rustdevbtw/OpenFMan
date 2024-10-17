package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Auth middleware for password protection using Bearer token with custom SHA-256 hashing
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Get the plain password from config
		password := viper.GetString("auth.password")

		// Generate the expected token based on the custom hashing logic
		expectedToken := generateToken(password)

		// Compare the provided token with the expected token (deterministic comparison)
		if token != expectedToken {
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

// generateToken generates the SHA-256-hashed token for comparison based on the custom logic
func generateToken(password string) string {
	// First apply SHA-256 twice and concatenate "_pass" to the result
	h1 := hashSHA256(password)
	h2 := hashSHA256(h1)
	s := "_pass" + h2

	// Hash the final result deterministically using SHA-256 again
	return hashSHA256(s)
}

// hashSHA256 returns the SHA-256 hash of a given string
func hashSHA256(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	return hex.EncodeToString(hash.Sum(nil))
}
