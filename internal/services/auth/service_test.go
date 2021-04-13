package auth

import (
	"net/http"
	"testing"
)

func TestGetToken(t *testing.T) {
	req, err := http.NewRequest("POST", "/random", nil)
	if err != nil {
		t.Error(err.Error())
	}
	t.Run("GetToken returns an empty string if a request without an authorization header is provided", func(t *testing.T) {
		token := GetToken(req)
		if token != "" {
			t.Errorf("expected an empty string, got: %s", token)
		}
	})
	t.Run("GetToken returns an empty string if a request with a wrong authorization header is provided", func(t *testing.T) {
		req.Header.Set("Authorization", "asas")
		token := GetToken(req)
		if token != "" {
			t.Errorf("expected an empty string, got: %s", token)
		}
	})
	t.Run("GetToken returns an empty string if a request has bearer authorization header without any token provided", func(t *testing.T) {
		req.Header.Set("Authorization", "Bearer")
		token := GetToken(req)
		if token != "" {
			t.Errorf("expected an empty string, got: %s", token)
		}
	})
	t.Run("GetToken should return a token if authorization header is valid", func(t *testing.T) {
		passedToken := "asd"
		req.Header.Set("Authorization", "Bearer "+passedToken)
		token := GetToken(req)
		if token != passedToken {
			t.Errorf("expected %s, got: %s", passedToken, token)
		}
	})
}
