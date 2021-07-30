package auth

import (
	"database/sql"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/dchest/uniuri"

	"github.com/adjsky/fetchapp_server/internal/models/user"
	"github.com/adjsky/fetchapp_server/pkg/helpers"

	"github.com/adjsky/fetchapp_server/config"
	"github.com/adjsky/fetchapp_server/internal/services"
	"github.com/gin-gonic/gin"
)

var emailRegex *regexp.Regexp

const (
	restoreSessionDuration = time.Minute * 15
	refreshPeriod          = time.Minute * 5
)

func init() {
	emailRegex = regexp.MustCompile(`^\S+@\S+$`)
}

type restoreSession struct {
	email     string
	createdAt time.Time
}

type authService struct {
	config          *config.Config
	database        *sql.DB
	restoreSessions map[string]restoreSession
	restoreMutex    sync.RWMutex
	userManager     *user.Manager
}

// NewService creates a new auth Service
func NewService(cfg *config.Config, db *sql.DB) services.Service {
	s := &authService{
		config:          cfg,
		database:        db,
		restoreSessions: make(map[string]restoreSession),
		userManager:     user.NewManager(db),
	}
	go func() {
		for {
			time.Sleep(refreshPeriod)
			s.CheckExpire()
		}
	}()
	return s
}

// Register the auth service
func (serv *authService) Register(r *gin.RouterGroup) {
	r.POST("/login", serv.handleLogin)
	r.POST("/signup", serv.handleSignup)
	r.PUT("/restore", serv.handleRestore)
	r.POST("/restore/valid", serv.handleRestoreValid)
	r.POST("/valid", serv.handleValid)
}

// Close does clean up actions on the service
func (serv *authService) Close() {
	//
}

// CheckExpire checks and deletes outdated restore tokens
func (serv *authService) CheckExpire() {
	for k, v := range serv.restoreSessions {
		timePassed := time.Since(v.createdAt)
		if timePassed.Seconds() >= restoreSessionDuration.Seconds() {
			delete(serv.restoreSessions, k)
		}
	}
}

func (serv *authService) handleLogin(c *gin.Context) {
	var reqData loginRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		helpers.ResponseInvalidBody(c)
		return
	}
	model, err := serv.userManager.MatchPassword(reqData.Email, reqData.Password)
	if err != nil {
		code := http.StatusUnauthorized
		c.JSON(code, gin.H{
			"code":    code,
			"message": err.Error(),
		})
		return
	}
	token, err := model.GetAuthToken(serv.config.SecretKey)
	if err != nil {
		code := http.StatusInternalServerError
		c.JSON(code, gin.H{
			"code":    code,
			"message": err.Error(),
		})
		return
	}
	code := http.StatusOK
	c.JSON(code, gin.H{
		"code":  code,
		"token": token,
	})
}

func (serv *authService) handleSignup(c *gin.Context) {
	var reqData signupRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		helpers.ResponseInvalidBody(c)
		return
	}
	matched := emailRegex.Match([]byte(reqData.Email))
	if !matched {
		code := http.StatusBadRequest
		c.JSON(code, gin.H{
			"code":    code,
			"message": "invalid email address",
		})
		return
	}
	model, err := serv.userManager.Create(reqData.Email, reqData.Password)
	if err != nil {
		var code int
		if err == user.ErrInternal {
			code = http.StatusInternalServerError
		} else if err == user.ErrEmailRegistered {
			code = http.StatusConflict
		}
		c.JSON(code, gin.H{
			"code":    code,
			"message": err.Error(),
		})
		return
	}
	token, err := model.GetAuthToken(serv.config.SecretKey)
	if err != nil {
		code := http.StatusInternalServerError
		c.JSON(code, gin.H{
			"code":    code,
			"message": err.Error(),
		})
		return
	}
	code := http.StatusOK
	c.JSON(code, gin.H{
		"code":  code,
		"token": token,
	})
}

func (serv *authService) handleRestore(c *gin.Context) {
	if CheckAuthorized(c) {
		serv.handleRestoreAuth(c)
	} else {
		serv.handleRestoreNotAuth(c)
	}
}

func (serv *authService) handleRestoreAuth(c *gin.Context) {
	var reqData restoreRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		helpers.ResponseInvalidBody(c)
		return
	}
	model, err := serv.userManager.GetModelFromToken(GetToken(c), serv.config.SecretKey)
	if err != nil {
		code := http.StatusBadRequest
		c.JSON(code, gin.H{
			"code":    code,
			"message": "invalid auth token provided",
		})
		return
	}
	err = serv.userManager.ChangePassword(model.Email, reqData.OldPassword, reqData.NewPassword)
	if err != nil {
		var code int
		if err == user.ErrNotMatched {
			code = http.StatusUnauthorized
		} else if err == user.ErrInternal {
			code = http.StatusInternalServerError
		}
		c.JSON(code, gin.H{
			"code":    code,
			"message": err.Error(),
		})
		return
	}
	code := http.StatusOK
	c.JSON(code, gin.H{
		"code": code,
	})
}

func (serv *authService) handleRestoreNotAuth(c *gin.Context) {
	var reqData restoreRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		helpers.ResponseInvalidBody(c)
		return
	}
	if reqData.Code == "" {
		isRegistered := serv.userManager.IsEmailRegistered(reqData.Email)
		if !isRegistered {
			code := http.StatusBadRequest
			c.JSON(code, gin.H{
				"code":    code,
				"message": user.ErrNoUser.Error(),
			})
			return
		}
		code := uniuri.NewLen(8)
		serv.restoreMutex.Lock()
		for k, v := range serv.restoreSessions {
			if v.email == reqData.Email {
				delete(serv.restoreSessions, k)
			}
		}
		serv.restoreSessions[code] = restoreSession{
			email:     reqData.Email,
			createdAt: time.Now(),
		}
		serv.restoreMutex.Unlock()
		statusCode := http.StatusAccepted
		c.JSON(statusCode, gin.H{
			"code": statusCode,
		})
		go func() {
			err := helpers.SendEmail(&serv.config.SMTP,
				[]string{reqData.Email},
				[]byte("Subject: Restore account\n"+code))
			if err != nil {
				log.Fatal(err)
			}
		}()
	} else {
		if reqData.NewPassword == "" || reqData.OldPassword == "" {
			helpers.ResponseInvalidBody(c)
			return
		}
		restoreSession, ok := serv.restoreSessions[reqData.Code]
		if !ok || restoreSession.email != reqData.Email {
			code := http.StatusBadRequest
			c.JSON(code, gin.H{
				"code":    code,
				"message": "invalid code provided",
			})
			return
		}
		err := serv.userManager.ChangePassword(reqData.Email, reqData.OldPassword, reqData.NewPassword)
		if err != nil {
			var code int
			if err == user.ErrNotMatched {
				code = http.StatusUnauthorized
			} else if err == user.ErrInternal {
				code = http.StatusInternalServerError
			}
			c.JSON(code, gin.H{
				"code":    code,
				"message": err.Error(),
			})
			return
		}
		delete(serv.restoreSessions, reqData.Code)
		code := http.StatusOK
		c.JSON(code, gin.H{
			"code": code,
		})
	}
}

func (serv *authService) handleValid(c *gin.Context) {
	var reqData validRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		helpers.ResponseInvalidBody(c)
		return
	}
	_, err := serv.userManager.GetModelFromToken(reqData.Token, serv.config.SecretKey)
	code := http.StatusOK
	c.JSON(code, gin.H{
		"code":  code,
		"valid": err == nil,
	})
}

func (serv *authService) handleRestoreValid(c *gin.Context) {
	var reqData restoreValidRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		helpers.ResponseInvalidBody(c)
		return
	}
	_, ok := serv.restoreSessions[reqData.Code]
	code := http.StatusOK
	c.JSON(code, gin.H{
		"code":  code,
		"valid": ok,
	})
}

// CheckAuthorized checks whether a given request has a bearer token
func CheckAuthorized(c *gin.Context) bool {
	authHeader := c.GetHeader("Authorization")
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
func GetToken(c *gin.Context) string {
	if CheckAuthorized(c) {
		return strings.Split(c.GetHeader("Authorization"), " ")[1]
	}
	return ""
}
