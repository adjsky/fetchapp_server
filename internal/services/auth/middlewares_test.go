package auth

import (
	"net/http"
	"net/http/httptest"
	"server/config"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	authService := NewService(config.Get(), nil)
	req, err := http.NewRequest("POST", "/random", nil)
	if err != nil {
		t.Error(err.Error())
	}
	handler := authService.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Run("Request with a null authorization header returns 401 status code", func(t *testing.T) {
		writer := httptest.NewRecorder()
		handler.ServeHTTP(writer, req)
		if writer.Code != http.StatusUnauthorized {
			t.Errorf("expected status code: %v, got: %v", http.StatusUnauthorized, writer.Code)
		}
	})
	t.Run("Request with a wrong authorization header returns 401 status code", func(t *testing.T) {
		writer := httptest.NewRecorder()
		req.Header.Set("Authorization", "Basic")
		handler.ServeHTTP(writer, req)
		if writer.Code != http.StatusUnauthorized {
			t.Errorf("expected status code: %v, got: %v", http.StatusUnauthorized, writer.Code)
		}
	})
	t.Run("Request with an invalid token returns 401 status code", func(t *testing.T) {
		writer := httptest.NewRecorder()
		req.Header.Set("Authorization", "Bearer "+"asd")
		handler.ServeHTTP(writer, req)
		if writer.Code != http.StatusUnauthorized {
			t.Errorf("expected status code: %v, got: %v", http.StatusUnauthorized, writer.Code)
		}
	})
	t.Run("Middleware should pass a request with a valid token", func(t *testing.T) {
		writer := httptest.NewRecorder()
		claims := GenerateClaims("loh@mail.ru")
		token, _ := GenerateTokenString(claims, authService.config.SecretKey)
		req.Header.Set("Authorization", "Bearer "+token)
		handler.ServeHTTP(writer, req)
		if writer.Code != http.StatusOK {
			t.Errorf("expected status code: %v, got: %v", http.StatusOK, writer.Code)
		}
	})
}
