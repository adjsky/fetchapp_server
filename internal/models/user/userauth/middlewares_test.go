package userauth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adjsky/fetchapp_server/config"
	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	cfg, err := config.Get()
	if err != nil {
		t.Fatal(err)
	}

	handler := Middleware(cfg.SecretKey)
	t.Run("Request with a null authorization header returns 401 status code",
		func(t *testing.T) {
			writer := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(writer)
			req, _ := http.NewRequest("POST", "/asdasd", nil)
			ctx.Request = req
			handler(ctx)
			if writer.Code != http.StatusUnauthorized {
				t.Errorf("expected status code: %v, got: %v", http.StatusUnauthorized, writer.Code)
			}
		})
	t.Run("Request with a wrong authorization header returns 401 status code",
		func(t *testing.T) {
			writer := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(writer)
			req, _ := http.NewRequest("POST", "/asdasd", nil)
			ctx.Request = req
			handler(ctx)
			if writer.Code != http.StatusUnauthorized {
				t.Errorf("expected status code: %v, got: %v", http.StatusUnauthorized, writer.Code)
			}
		})
	t.Run("Request with an invalid token returns 401 status code",
		func(t *testing.T) {
			writer := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(writer)
			req, _ := http.NewRequest("POST", "/asdasd", nil)
			req.Header.Set("Authorization", "Bearer "+"asd")
			ctx.Request = req
			handler(ctx)
			if writer.Code != http.StatusUnauthorized {
				t.Errorf("expected status code: %v, got: %v", http.StatusUnauthorized, writer.Code)
			}
		})
	t.Run("Middleware should pass a request with a valid token",
		func(t *testing.T) {
			claims := GenerateClaims("loh@mail.ru")
			token, _ := GenerateToken(claims, cfg.SecretKey)
			writer := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(writer)
			req, _ := http.NewRequest("POST", "/asdasd", nil)
			req.Header.Set("Authorization", "Bearer "+token)
			ctx.Request = req
			handler(ctx)
			if writer.Code != http.StatusOK {
				t.Errorf("expected status code: %v, got: %v", http.StatusOK, writer.Code)
			}
		})
}
