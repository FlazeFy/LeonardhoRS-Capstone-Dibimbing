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

func ValidateToken(tokenString string) (uuid.UUID, error) {
	// Token Parse
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return config.GetJWTSecret(), nil
	})
	if err != nil || !token.Valid {
		return uuid.UUID{}, errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.UUID{}, errors.New("invalid token claims")
	}

	// Check exp
	expFloat, ok := claims["exp"].(float64)
	if !ok {
		return uuid.UUID{}, errors.New("missing exp in token")
	}

	expTime := time.Unix(int64(expFloat), 0)
	if time.Now().After(expTime) {
		return uuid.UUID{}, errors.New("token expired")
	}

	userIdStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.UUID{}, errors.New("user_id is not a string")
	}

	userID, err := uuid.Parse(userIdStr)
	if err != nil {
		return uuid.UUID{}, errors.New("invalid user_id format")
	}

	return userID, nil
}

func AuthMiddleware(redisClient *redis.Client, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization header"})
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

		// Validate JWT Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unexpected signing method"})
			}
			return config.GetJWTSecret(), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid or expired token"})
			return
		}

		// Claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
			return
		}

		// Extract userID
		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "user ID not found"})
			return
		}

		// Extract role
		role, ok := claims["role"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "role not found"})
			return
		}

		// Check If Role Is Allowed
		isAllowed := false
		for _, r := range allowedRoles {
			if role == r {
				isAllowed = true
				break
			}
		}
		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "access forbidden for this role"})
			return
		}

		// Set Context
		c.Set("userID", userID)
		c.Set("role", role)

		c.Next()
	}
}
