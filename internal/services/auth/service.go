package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"server/config"
	"server/pkg/handlers"
	"server/pkg/helpers"
	"server/pkg/middlewares"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"
)

var emailRegex *regexp.Regexp

const restoreSessionDuration = time.Minute * 15

func init() {
	emailRegex = regexp.MustCompile(`^\S+@\S+$`)
}

type restoreSession struct {
	email     string
	createdAt time.Time
}

type service struct {
	config          *config.Config
	database        *sql.DB
	restoreSessions map[string]restoreSession
	restoreMutex    sync.RWMutex
}

// NewService creates a new auth service
func NewService(cfg *config.Config, db *sql.DB) *service {
	s := service{
		config:          cfg,
		database:        db,
		restoreSessions: make(map[string]restoreSession),
	}
	return &s
}

// Register auth service
func (serv *service) Register(r *mux.Router) {
	appJsonMiddleware := middlewares.ContentTypeValidator("application/json")
	r.Handle("/login", appJsonMiddleware(http.HandlerFunc(serv.handleLogin))).Methods("POST")
	r.Handle("/signup", appJsonMiddleware(http.HandlerFunc(serv.handleSignup))).Methods("POST")
	r.Handle("/restore", appJsonMiddleware(http.HandlerFunc(serv.handleRestore))).Methods("PUT")
	r.Handle("/restore/valid", appJsonMiddleware(http.HandlerFunc(serv.handleRestoreValid))).Methods("POST")
	r.Handle("/valid", appJsonMiddleware(http.HandlerFunc(serv.handleValid))).Methods("POST")
}

// CheckExpire checks and deletes outdated restore tokens
func (serv *service) CheckExpire() {
	for k, v := range serv.restoreSessions {
		timePassed := time.Since(v.createdAt)
		if timePassed.Seconds() >= restoreSessionDuration.Seconds() {
			delete(serv.restoreSessions, k)
		}
	}
}

func (serv *service) handleLogin(w http.ResponseWriter, req *http.Request) {
	data, _ := io.ReadAll(req.Body)
	var reqData loginRequest
	err := json.Unmarshal(data, &reqData)
	if err != nil {
		handlers.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if reqData.Email == "" || reqData.Password == "" {
		handlers.RespondError(w, http.StatusBadRequest, "no password or email provided")
		return
	}
	var password string
	row := serv.database.QueryRow("SELECT password FROM Users WHERE email = ?", reqData.Email)
	err = row.Scan(&password)
	if err != nil {
		handlers.RespondError(w, http.StatusUnauthorized, "no user registered with this email")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(password), []byte(reqData.Password)) != nil {
		handlers.RespondError(w, http.StatusUnauthorized, "wrong email/password pair")
		return
	}
	claims := GenerateClaims(reqData.Email)
	token, _ := GenerateTokenString(claims, serv.config.SecretKey)
	res := loginResponse{
		Code:  http.StatusOK,
		Token: token,
	}
	handlers.Respond(w, &res, res.Code)
}

func (serv *service) handleSignup(w http.ResponseWriter, req *http.Request) {
	data, _ := io.ReadAll(req.Body)
	var reqData signupRequest
	err := json.Unmarshal(data, &reqData)
	if err != nil {
		handlers.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if reqData.Email == "" || reqData.Password == "" {
		handlers.RespondError(w, http.StatusBadRequest, "no password or email provided")
		return
	}
	matched := emailRegex.Match([]byte(reqData.Email))
	if !matched {
		handlers.RespondError(w, http.StatusBadRequest, "invalid email address")
		return
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(reqData.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("hash generating error in signup: ", err)
	}
	_, err = serv.database.Exec("INSERT INTO Users (email, password) VALUES (?, ?)", reqData.Email, hashPassword)
	if err != nil {
		handlers.RespondError(w, http.StatusConflict, "this email is registered")
		return
	}
	claims := GenerateClaims(reqData.Email)
	token, err := GenerateTokenString(claims, serv.config.SecretKey)
	if err != nil {
		log.Println(err)
	}
	res := signupResponse{
		Code:  http.StatusOK,
		Token: token,
	}
	handlers.Respond(w, &res, res.Code)
}

func (serv *service) handleRestore(w http.ResponseWriter, req *http.Request) {
	if CheckAuthorized(req) {
		serv.handleRestoreAuth(w, req)
	} else {
		serv.handleRestoreNonAuth(w, req)
	}
}

func (serv *service) handleRestoreAuth(w http.ResponseWriter, req *http.Request) {
	data, _ := io.ReadAll(req.Body)
	var reqData restoreRequest
	err := json.Unmarshal(data, &reqData)
	if err != nil {
		handlers.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if reqData.OldPassword == "" || reqData.NewPassword == "" {
		handlers.RespondError(w, http.StatusBadRequest, "no old or new password provided")
		return
	}
	userClaims, err := GetClaims(GetToken(req), serv.config.SecretKey)
	if err != nil {
		handlers.RespondError(w, http.StatusBadRequest, "invalid auth token provided")
		return
	}
	var userID int
	var userPassword string
	row := serv.database.QueryRow("SELECT ID, password FROM Users WHERE email = ?", userClaims.Email)
	_ = row.Scan(&userID, &userPassword)
	if bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(reqData.OldPassword)) != nil {
		handlers.RespondError(w, http.StatusUnauthorized, "old password doesn't correspond to account password")
		return
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(reqData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("hash generating error in restore: ", err)
	}
	_, _ = serv.database.Exec("UPDATE Users SET password = ? WHERE ID = ?", hashPassword, userID)
	handlers.Respond(w, restoreResponse{Code: http.StatusOK}, http.StatusOK)
}

func (serv *service) handleRestoreNonAuth(w http.ResponseWriter, req *http.Request) {
	data, _ := io.ReadAll(req.Body)
	var reqData restoreRequest
	err := json.Unmarshal(data, &reqData)
	if err != nil {
		handlers.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if reqData.Code == "" {
		if reqData.Email == "" {
			handlers.RespondError(w, http.StatusBadRequest, "no email provided")
			return
		}
		row := serv.database.QueryRow("SELECT ID FROM Users WHERE email = ?", reqData.Email)
		err := row.Scan()
		if err == sql.ErrNoRows {
			handlers.RespondError(w, http.StatusBadRequest, "no user with provided email registered")
			return
		}
		code := uniuri.NewLen(8)
		serv.restoreMutex.Lock()
		defer serv.restoreMutex.Unlock()
		for k, v := range serv.restoreSessions {
			if v.email == reqData.Email {
				delete(serv.restoreSessions, k)
			}
		}
		serv.restoreSessions[code] = restoreSession{
			email:     reqData.Email,
			createdAt: time.Now(),
		}
		res := restoreResponse{
			Code: http.StatusAccepted,
		}
		handlers.Respond(w, &res, res.Code)
		go func() {
			err := helpers.SendEmail(&serv.config.Smtp,
				[]string{reqData.Email},
				[]byte("Subject: Restore account\n"+code))
			if err != nil {
				fmt.Println(err)
			}
		}()
	} else {
		if reqData.NewPassword == "" || reqData.OldPassword == "" {
			handlers.RespondError(w, http.StatusBadRequest, "no new or old password provided")
			return
		}
		restoreSession, ok := serv.restoreSessions[reqData.Code]
		if !ok {
			handlers.RespondError(w, http.StatusBadRequest, "invalid token provided")
			return
		}
		var userID int
		var userPassword string
		row := serv.database.QueryRow("SELECT ID, password FROM Users WHERE email = ?", restoreSession.email)
		_ = row.Scan(&userID, &userPassword)
		if bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(reqData.OldPassword)) != nil {
			handlers.RespondError(w, http.StatusUnauthorized, "old password doesn't correspond to account password")
			return
		}
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(reqData.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("hash generating error in restore: ", err)
		}
		_, _ = serv.database.Exec("UPDATE Users SET password = ? WHERE ID = ?", hashPassword, userID)
		delete(serv.restoreSessions, reqData.Code)
		handlers.Respond(w, restoreResponse{Code: http.StatusOK}, http.StatusOK)
	}
}

func (serv *service) handleValid(w http.ResponseWriter, req *http.Request) {
	data, _ := io.ReadAll(req.Body)
	var reqData validRequest
	err := json.Unmarshal(data, &reqData)
	if err != nil {
		handlers.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if reqData.Token == "" {
		handlers.RespondError(w, http.StatusBadRequest, "no token provided")
		return
	}
	_, err = GetClaims(reqData.Token, serv.config.SecretKey)
	res := validResponse{
		Code:  http.StatusOK,
		Valid: err == nil,
	}
	handlers.Respond(w, &res, res.Code)
}

func (serv *service) handleRestoreValid(w http.ResponseWriter, req *http.Request) {
	data, _ := io.ReadAll(req.Body)
	var reqData restoreValidRequest
	err := json.Unmarshal(data, &reqData)
	if err != nil {
		handlers.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if reqData.Code == "" || reqData.Email == "" {
		handlers.RespondError(w, http.StatusBadRequest, "no code or email provided")
		return
	}
	_, ok := serv.restoreSessions[reqData.Code]
	res := validResponse{
		Code:  http.StatusOK,
		Valid: ok,
	}
	handlers.Respond(w, &res, res.Code)
}

// CheckAuthorized checks whether a given request has a bearer token and returns it
func CheckAuthorized(req *http.Request) bool {
	authHeader := req.Header.Get("Authorization")
	authData := strings.Split(authHeader, " ")
	if len(authData) == 0 {
		return false
	}
	if authData[0] != "Bearer" || len(authData) != 2 {
		return false
	}
	return true
}

// GetToken returns a token from a request or an empty string if an user is not authorized or has an invalid authorization header
func GetToken(req *http.Request) string {
	if CheckAuthorized(req) {
		return strings.Split(req.Header.Get("Authorization"), " ")[1]
	}
	return ""
}
