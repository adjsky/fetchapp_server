package auth

import (
	"context"
	"net/http"
	"server/pkg/handlers"
	"strings"
)

// ContextKey is a type that is used to reference data in a context associated with
type ContextKey int

const (
	// ClaimsID constant is used to reference claims in a request context
	ClaimsID ContextKey = iota + 1
)

// AuthMiddleware checks whether a user has JWT token
func (serv *service) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		authData := strings.Split(authHeader, " ")
		if len(authData) == 0 {
			handlers.RespondError(w, http.StatusUnauthorized, "no authorization header provided")
			return
		}
		if authData[0] != "Bearer" {
			handlers.RespondError(w, http.StatusUnauthorized, "wrong authorization method provided")
			return
		}
		if len(authData) != 2 {
			handlers.RespondError(w, http.StatusUnauthorized, "no token provided")
			return
		}
		claims, err := GetClaims(authData[1], serv.config.SecretKey)
		if err != nil {
			handlers.RespondError(w, http.StatusUnauthorized, "invalid auth token provided")
			return
		}
		req = req.WithContext(context.WithValue(req.Context(), ClaimsID, claims))
		next.ServeHTTP(w, req)
	})
}
