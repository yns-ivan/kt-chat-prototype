package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"ktchat/backend/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CORS middleware for handling Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Logger middleware for request logging
func Logger() gin.HandlerFunc {
	return gin.Logger()
}

// AuthMiddleware validates JWT tokens (both custom and Cognito)
func AuthMiddleware(jwtSecret string, cognitoAuth *auth.CognitoService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Try to validate as custom JWT token first
		if userInfo, err := validateCustomToken(tokenString, jwtSecret); err == nil {
			// Set user information in context
			c.Set("user_id", userInfo.ID)
			c.Set("username", userInfo.Username)
			c.Set("email", userInfo.Email)
			c.Next()
			return
		}

		// If custom token validation fails, try Cognito token
		if cognitoAuth != nil {
			if userInfo, err := cognitoAuth.ValidateToken(c.Request.Context(), tokenString); err == nil {
				// Set user information in context
				c.Set("user_id", userInfo.ID)
				c.Set("username", userInfo.Username)
				c.Set("email", userInfo.Email)
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
	}
}

// validateCustomToken validates a custom JWT token
func validateCustomToken(tokenString, jwtSecret string) (*auth.UserInfo, error) {
	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Validate issuer and audience
	if iss, ok := claims["iss"].(string); !ok || iss != "ktchat-backend" {
		return nil, fmt.Errorf("invalid token issuer")
	}

	if aud, ok := claims["aud"].(string); !ok || aud != "ktchat-frontend" {
		return nil, fmt.Errorf("invalid token audience")
	}

	// Extract user information
	userInfo := &auth.UserInfo{}

	if userID, ok := claims["user_id"].(string); ok {
		userInfo.ID = userID
	} else {
		return nil, fmt.Errorf("missing user_id in token")
	}

	if username, ok := claims["username"].(string); ok {
		userInfo.Username = username
	} else {
		return nil, fmt.Errorf("missing username in token")
	}

	if email, ok := claims["email"].(string); ok {
		userInfo.Email = email
	} else {
		return nil, fmt.Errorf("missing email in token")
	}

	return userInfo, nil
}

// RateLimit middleware for basic rate limiting
func RateLimit() gin.HandlerFunc {
	// Simple in-memory rate limiter (in production, use Redis)
	clients := make(map[string][]time.Time)
	
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		now := time.Now()
		
		// Clean old entries (older than 1 minute)
		if times, exists := clients[clientIP]; exists {
			var validTimes []time.Time
			for _, t := range times {
				if now.Sub(t) < time.Minute {
					validTimes = append(validTimes, t)
				}
			}
			clients[clientIP] = validTimes
		}
		
		// Check rate limit (100 requests per minute)
		if len(clients[clientIP]) >= 100 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		
		// Add current request
		clients[clientIP] = append(clients[clientIP], now)
		
		c.Next()
	}
} 