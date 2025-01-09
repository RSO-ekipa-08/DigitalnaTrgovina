package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header"})
            c.Abort()
            return
        }

        // Extract the token from the Authorization header
        bearerToken := strings.Split(authHeader, " ")
        if len(bearerToken) != 2 {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
            c.Abort()
            return
        }

        token := bearerToken[1]

        // Verify token here if needed
        // For Auth0, the token is already verified by Auth0

        // Add token claims to context if needed
        c.Set("token", token)
        c.Next()
    }
}
