package middlewares

import (
	"context"
	"fmt"
	"log"
	"mime"
	"net/http"
	"server/pkg/handlers"
	"strings"
)

// ContextKey is a type that is used to reference data in a context associated with
type ContextKey int

const (
	// BoundaryID constant is used to reference a boundary provided with a multipart request
	BoundaryID ContextKey = iota + 1
)

// Log middleware logs every incoming call to the server
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println(req.URL.Path, req.Method, req.UserAgent(), req.RemoteAddr)
		next.ServeHTTP(w, req)
	})
}

// ContentTypeValidator returns a middleware that checks whether a request header equals to required
func ContentTypeValidator(required string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			contentType, params, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))
			if contentType == "" {
				handlers.RespondError(w, http.StatusBadRequest, "no content-type header provided")
				return
			}
			if contentType != required {
				handlers.RespondError(w, http.StatusBadRequest, fmt.Sprintf("server accepts requests only with %s header", required))
				return
			}
			if strings.HasPrefix(contentType, "multipart") {
				boundary, ok := params["boundary"]
				if !ok {
					handlers.RespondError(w, http.StatusBadRequest, "no boundary provided")
					return
				}
				req = req.WithContext(context.WithValue(req.Context(), BoundaryID, boundary))
			}
			next.ServeHTTP(w, req)
		})
	}
}
