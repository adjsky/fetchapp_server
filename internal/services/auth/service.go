package auth

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"server/config"
	"server/pkg/helpers"
	"strings"
	"sync"
	"time"

	"github.com/dchest/uniuri"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

var emailRegex *regexp.Regexp

const (
	restoreSessionDuration = time.Minute * 15
)

func init() {
	emailRegex = regexp.MustCompile(`^\S+@\S+$`)
}

type RestoreSession struct {
	Email     string
	CreatedAt time.Time
}

type Service struct {
	Config          *config.Config
	Database        *sql.DB
	RestoreSessions map[string]RestoreSession
	RestoreMutex    sync.RWMutex
}

// NewService creates a new auth Service
func NewService(cfg *config.Config, db *sql.DB) *Service {
	s := Service{
		Config:          cfg,
		Database:        db,
		RestoreSessions: make(map[string]RestoreSession),
	}
	return &s
}

// Register auth Service
func (serv *Service) Register(r *gin.RouterGroup) {
	r.POST("/login", serv.handleLogin)
	r.POST("/signup", serv.handleSignup)
	r.PUT("/restore", serv.handleRestore)
	r.POST("/restore/valid", serv.handleRestoreValid)
	r.POST("/valid", serv.handleValid)
}

// CheckExpire checks and deletes outdated restore tokens
func (serv *Service) CheckExpire() {
	for k, v := range serv.RestoreSessions {
		timePassed := time.Since(v.CreatedAt)
		if timePassed.Seconds() >= restoreSessionDuration.Seconds() {
			delete(serv.RestoreSessions, k)
		}
	}
}

func (serv *Service) handleLogin(c *gin.Context) {
	var reqData loginRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		code := http.StatusBadRequest
		c.JSON(code, gin.H{
			"code":    code,
			"message": "invalid request body",
		})
		return
	}
	var password string
	row := serv.Database.QueryRow("SELECT password FROM Users WHERE email = $1", reqData.Email)
	if err := row.Scan(&password); err != nil {
		code := http.StatusUnauthorized
		c.JSON(code, gin.H{
			"code":    code,
			"message": "no user registered with this email",
		})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(password), []byte(reqData.Password)) != nil {
		code := http.StatusUnauthorized
		c.JSON(code, gin.H{
			"code":    code,
			"message": "wrong email/password pair",
		})
		return
	}
	claims := GenerateClaims(reqData.Email)
	token, _ := GenerateTokenString(claims, serv.Config.SecretKey)
	code := http.StatusOK
	c.JSON(code, gin.H{
		"code":  code,
		"token": token,
	})
}

func (serv *Service) handleSignup(c *gin.Context) {
	var reqData signupRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		code := http.StatusBadRequest
		c.JSON(code, gin.H{
			"code":    code,
			"message": "invalid request body",
		})
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
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(reqData.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("hash generating error in signup: ", err)
	}
	_, err = serv.Database.Exec("INSERT INTO Users (email, password) VALUES ($1, $2)", reqData.Email, hashPassword)
	if err != nil {
		code := http.StatusConflict
		c.JSON(code, gin.H{
			"code":    code,
			"message": "this email is registered",
		})
		return
	}
	claims := GenerateClaims(reqData.Email)
	token, err := GenerateTokenString(claims, serv.Config.SecretKey)
	if err != nil {
		log.Println(err)
	}
	code := http.StatusOK
	c.JSON(code, gin.H{
		"code":  code,
		"token": token,
	})
}

func (serv *Service) handleRestore(c *gin.Context) {
	if CheckAuthorized(c) {
		serv.handleRestoreAuth(c)
	} else {
		serv.handleRestoreNonAuth(c)
	}
}

func (serv *Service) handleRestoreAuth(c *gin.Context) {
	var reqData restoreRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		code := http.StatusBadRequest
		c.JSON(code, gin.H{
			"code":    code,
			"message": "invalid request body",
		})
		return
	}
	userClaims, err := GetClaims(GetToken(c), serv.Config.SecretKey)
	if err != nil {
		code := http.StatusBadRequest
		c.JSON(code, gin.H{
			"code":    code,
			"message": "invalid auth token provided",
		})
		return
	}
	var userID int
	var userPassword string
	row := serv.Database.QueryRow("SELECT ID, password FROM Users WHERE email = $1", userClaims.Email)
	_ = row.Scan(&userID, &userPassword)
	if bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(reqData.OldPassword)) != nil {
		code := http.StatusUnauthorized
		c.JSON(code, gin.H{
			"code":    code,
			"message": "old password doesn't correspond to account password",
		})
		return
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(reqData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("hash generating error in restore: ", err)
	}
	_, _ = serv.Database.Exec("UPDATE Users SET password = $1 WHERE ID = $2", hashPassword, userID)
	code := http.StatusOK
	c.JSON(code, gin.H{
		"code": code,
	})
}

func (serv *Service) handleRestoreNonAuth(c *gin.Context) {
	var reqData restoreRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		code := http.StatusBadRequest
		c.JSON(code, gin.H{
			"code":    code,
			"message": "invalid request body",
		})
		return
	}
	if reqData.Code == "" {
		if reqData.Email == "" {
			code := http.StatusBadRequest
			c.JSON(code, gin.H{
				"code":    code,
				"message": "no email provided",
			})
			return
		}
		row := serv.Database.QueryRow("SELECT ID FROM Users WHERE email = $1", reqData.Email)
		if err := row.Scan(); err == sql.ErrNoRows {
			code := http.StatusBadRequest
			c.JSON(code, gin.H{
				"code":    code,
				"message": "no user with provided email registered",
			})
			return
		}
		code := uniuri.NewLen(8)
		serv.RestoreMutex.Lock()
		defer serv.RestoreMutex.Unlock()
		for k, v := range serv.RestoreSessions {
			if v.Email == reqData.Email {
				delete(serv.RestoreSessions, k)
			}
		}
		serv.RestoreSessions[code] = RestoreSession{
			Email:     reqData.Email,
			CreatedAt: time.Now(),
		}
		statusCode := http.StatusAccepted
		c.JSON(statusCode, gin.H{
			"code": statusCode,
		})
		go func() {
			err := helpers.SendEmail(&serv.Config.SMTP,
				[]string{reqData.Email},
				[]byte("Subject: Restore account\n"+code))
			if err != nil {
				fmt.Println(err)
			}
		}()
	} else {
		if reqData.NewPassword == "" || reqData.OldPassword == "" {
			code := http.StatusBadRequest
			c.JSON(code, gin.H{
				"code":    code,
				"message": "no new or old password provided",
			})
			return
		}
		restoreSession, ok := serv.RestoreSessions[reqData.Code]
		if !ok {
			code := http.StatusBadRequest
			c.JSON(code, gin.H{
				"code":    code,
				"message": "invalid token provided",
			})
			return
		}
		var userID int
		var userPassword string
		row := serv.Database.QueryRow("SELECT ID, password FROM Users WHERE email = $1", restoreSession.Email)
		_ = row.Scan(&userID, &userPassword)
		if bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(reqData.OldPassword)) != nil {
			code := http.StatusUnauthorized
			c.JSON(code, gin.H{
				"code":    code,
				"message": "old password doesn't correspond to account password",
			})
			return
		}
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(reqData.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("hash generating error in restore: ", err)
		}
		_, _ = serv.Database.Exec("UPDATE Users SET password = $1 WHERE ID = $2", hashPassword, userID)
		delete(serv.RestoreSessions, reqData.Code)
		code := http.StatusOK
		c.JSON(code, gin.H{
			"code": code,
		})
	}
}

func (serv *Service) handleValid(c *gin.Context) {
	var reqData validRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		code := http.StatusBadRequest
		c.JSON(code, gin.H{
			"code":    code,
			"message": "invalid request body",
		})
		return
	}
	_, err := GetClaims(reqData.Token, serv.Config.SecretKey)
	code := http.StatusOK
	c.JSON(code, gin.H{
		"code":  code,
		"valid": err == nil,
	})
}

func (serv *Service) handleRestoreValid(c *gin.Context) {
	var reqData restoreValidRequest
	if err := c.ShouldBindJSON(&reqData); err != nil {
		code := http.StatusBadRequest
		c.JSON(code, gin.H{
			"code":    code,
			"message": "invalid request body",
		})
		return
	}
	_, ok := serv.RestoreSessions[reqData.Code]
	code := http.StatusOK
	c.JSON(code, gin.H{
		"code":  code,
		"valid": ok,
	})
}

// CheckAuthorized checks whether a given request has a bearer token and returns it
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
