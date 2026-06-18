package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const ClaimsKey = "claims"

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			abortUnauthorized(c, "missing or invalid authorization header")
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			abortUnauthorized(c, "invalid token")
			return
		}

		c.Set(ClaimsKey, token.Claims)
		c.Next()
	}
}

func abortUnauthorized(c *gin.Context, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"error": gin.H{
			"code":    "UNAUTHORIZED",
			"message": message,
		},
	})
}
