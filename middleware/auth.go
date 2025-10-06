package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/aprianimmanuel/rangkaiedu-backend/config"
)

// UserContext represents the user information stored in the context
type UserContext struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
}

// AuthRequired is a middleware that validates JWT tokens and extracts user info
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header has the correct format "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			// Return the secret key for validation
			return []byte(config.Load().JWTSecret), nil
		})

		if err != nil {
			// Debug logging
			fmt.Printf("Token parsing error: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "debug": err.Error()})
			c.Abort()
			return
		}

		// Check if the token is valid
		if !token.Valid {
			fmt.Printf("Token is not valid\n")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "debug": "Token not valid"})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			fmt.Printf("Failed to extract claims from token\n")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Validate issuer
		if issuer, ok := claims["iss"].(string); !ok || issuer != "rangkai-edu-backend" {
			fmt.Printf("Invalid issuer: %v, expected: rangkai-edu-backend\n", issuer)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token issuer", "debug": fmt.Sprintf("Issuer: %v", issuer)})
			c.Abort()
			return
		}
		fmt.Printf("Token issuer validated: %v\n", claims["iss"])

		// Validate audience
		if audience, ok := claims["aud"].(string); !ok || audience != "rangkai-edu-frontend" {
			fmt.Printf("Invalid audience: %v, expected: rangkai-edu-frontend\n", audience)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token audience", "debug": fmt.Sprintf("Audience: %v", audience)})
			c.Abort()
			return
		}

		// Create user context from claims
		user := UserContext{
			ID:    getStringClaim(claims, "sub"),
			Email: getStringClaim(claims, "email"),
			Phone: getStringClaim(claims, "phone"),
			Role:  getStringClaim(claims, "role"),
		}

		// Store user context in Gin context
		c.Set("user", user)

		// Continue with the next handler
		c.Next()
	}
}

// getStringClaim safely extracts a string claim from jwt.MapClaims
func getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}