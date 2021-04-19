package auth

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetToken(t *testing.T) {
	t.Run("GetToken returns an empty string if a request without an authorization header is provided",
		func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(nil)
			req, _ := http.NewRequest("POST", "/random", nil)
			ctx.Request = req
			token := GetToken(ctx)
			if token != "" {
				t.Errorf("expected an empty string, got: %s", token)
			}
		})
	t.Run("GetToken returns an empty string if a request with a wrong authorization header is provided",
		func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(nil)
			req, _ := http.NewRequest("POST", "/random", nil)
			req.Header.Set("Authorization", "asas")
			ctx.Request = req
			token := GetToken(ctx)
			if token != "" {
				t.Errorf("expected an empty string, got: %s", token)
			}
		})
	t.Run("GetToken returns an empty string if a request has bearer authorization header without any token provided",
		func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(nil)
			req, _ := http.NewRequest("POST", "/random", nil)
			req.Header.Set("Authorization", "Bearer")
			ctx.Request = req
			token := GetToken(ctx)
			if token != "" {
				t.Errorf("expected an empty string, got: %s", token)
			}
		})
	t.Run("GetToken should return a token if authorization header is valid",
		func(t *testing.T) {
			passedToken := "asd"
			ctx, _ := gin.CreateTestContext(nil)
			req, _ := http.NewRequest("POST", "/random", nil)
			req.Header.Set("Authorization", "Bearer "+passedToken)
			ctx.Request = req
			token := GetToken(ctx)
			if token != passedToken {
				t.Errorf("expected %s, got: %s", passedToken, token)
			}
		})
}
