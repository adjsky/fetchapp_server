package auth

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/config"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	// initialize a driver
	_ "github.com/lib/pq"
)

func TestHandleLogin(t *testing.T) {
	cfg, err := config.Get()
	if err != nil {
		t.Fatal(err)
	}
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		t.Fatal(err)
	}
	userEmail := "testingemail@mail.ru"
	userPassword := "asdasd"
	userPasswordHash, _ := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	_, _ = db.Exec("INSERT INTO Users(email, password) VALUES($1, $2)", userEmail, userPasswordHash)
	service := authService{config: cfg, database: db}
	t.Run("handleLogin should return 400 status code if an invalid body is passed", func(t *testing.T) {
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "auth/login", nil)
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = request
		service.handleLogin(ctx)
		if writer.Code != http.StatusBadRequest {
			t.Errorf("expected: 400, got: %v", writer.Code)
		}
	})
	t.Run("handleLogin should return 401 status code if an unregistered email is passed", func(t *testing.T) {
		reqData, err := json.Marshal(loginRequest{
			Email:    "unregistered@gmail.com",
			Password: "123",
		})
		if err != nil {
			t.Error(err)
			return
		}
		body := bytes.NewReader(reqData)
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "auth/login", body)
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = request
		service.handleLogin(ctx)
		if writer.Code != 401 {
			t.Errorf("expected: 401, got: %v", writer.Code)
		}
	})
	t.Run("handleLogin should return 401 status code if an invalid email/password pair is passed", func(t *testing.T) {
		reqData, err := json.Marshal(loginRequest{
			Email:    userEmail,
			Password: "1",
		})
		if err != nil {
			t.Error(err)
			return
		}
		body := bytes.NewReader(reqData)
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "auth/login", body)
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = request
		service.handleLogin(ctx)
		if writer.Code != 401 {
			t.Errorf("expected: 401, got: %v", writer.Code)
		}
	})
	t.Run("handleLogin should return 200 status code if a valid email/password pair is passed", func(t *testing.T) {
		reqData, err := json.Marshal(loginRequest{
			Email:    userEmail,
			Password: userPassword,
		})
		if err != nil {
			t.Error(err)
			return
		}
		body := bytes.NewReader(reqData)
		writer := httptest.NewRecorder()
		request, _ := http.NewRequest("POST", "auth/login", body)
		ctx, _ := gin.CreateTestContext(writer)
		ctx.Request = request
		service.handleLogin(ctx)
		if writer.Code != 200 {
			t.Errorf("expected: 200, got: %v", writer.Code)
		}
	})
}

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
