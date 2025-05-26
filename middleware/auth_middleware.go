package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"pelita/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	// Token Parse
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.GetJWTSecret(), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Check exp
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("missing exp in token")
	}

	expTime := time.Unix(int64(expFloat), 0)
	if time.Now().After(expTime) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func AuthMiddleware(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing authorization header"})
			return
		}
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization header"})
			return
		}

		// Check If Token Is Blacklisted In Redis
		val, err := redisClient.Get(context.Background(), tokenString).Result()
		if err == nil && val == "blacklisted" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token already expired"})
			return
		}

		// Validate JWT token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if ok {
			userID, err := uuid.Parse(userIDStr)
			if err == nil {
				c.Set("userID", userID)
			}
		}

		c.Next()
	}
}
