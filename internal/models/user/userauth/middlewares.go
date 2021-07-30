package userauth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// ClaimsKey constant is used to reference claims in a request context
	ClaimsKey = "claims"
)

// Middleware checks whether a user has JWT token
func Middleware(secretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		authData := strings.Split(authHeader, " ")
		if len(authData) == 0 {
			code := http.StatusUnauthorized
			c.AbortWithStatusJSON(code, gin.H{
				"code":    code,
				"message": "no authorization header provided",
			})
			return
		}
		if authData[0] != "Bearer" {
			code := http.StatusUnauthorized
			c.AbortWithStatusJSON(code, gin.H{
				"code":    code,
				"message": "wrong authorization method provided",
			})
			return
		}
		if len(authData) != 2 {
			code := http.StatusUnauthorized
			c.AbortWithStatusJSON(code, gin.H{
				"code":    code,
				"message": "no token provided",
			})
			return
		}
		claims, err := GetClaims(authData[1], secretKey)
		if err != nil {
			code := http.StatusUnauthorized
			c.AbortWithStatusJSON(code, gin.H{
				"code":    code,
				"message": "invalid auth token provided",
			})
			return
		}
		c.Set(ClaimsKey, claims)
	}
}
