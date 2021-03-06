package user

import (
	"database/sql"
	"log"

	"github.com/adjsky/fetchapp_server/internal/models/user/userauth"

	"golang.org/x/crypto/bcrypt"
)

// Manager manages user models
type Manager struct {
	Database *sql.DB
}

// NewManager returns an user model manager
func NewManager(db *sql.DB) *Manager {
	return &Manager{
		Database: db,
	}
}

// Create creates a new user and returns a model
func (manager *Manager) Create(email, password string) (*Model, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrInternal
	}
	_, err = manager.Database.Exec("INSERT INTO Users (email, password) VALUES ($1, $2)", email, hashedPassword)
	if err != nil {
		log.Println("manager.Create error: " + err.Error())
		return nil, ErrEmailRegistered
	}
	return &Model{
		Email: email,
	}, nil
}

// MatchPassword checks whether the provided password matches and returns an user model
func (manager *Manager) MatchPassword(email, password string) (*Model, error) {
	var hashedPassword string
	row := manager.Database.QueryRow("SELECT password FROM Users WHERE email = $1", email)
	if err := row.Scan(&hashedPassword); err != nil {
		log.Println("manager.MatchPassword error: " + err.Error())
		if err == sql.ErrNoRows {
			return nil, ErrNoUser
		}
		return nil, ErrInternal
	}
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
		return nil, ErrNotMatched
	}
	return &Model{
		Email: email,
	}, nil
}

// ChangePassword changes an user password
func (manager *Manager) ChangePassword(email, oldPassword, newPassword string) error {
	var hashedPassword string
	row := manager.Database.QueryRow("SELECT password FROM Users WHERE email = $1", email)
	if err := row.Scan(&hashedPassword); err != nil {
		log.Println("manager.ChangePassword error: " + err.Error())
		if err == sql.ErrNoRows {
			return ErrNoUser
		}
		return ErrInternal
	}
	if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(oldPassword)) != nil {
		return ErrNotMatched
	}
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return ErrInternal
	}
	_, err = manager.Database.Exec("UPDATE Users SET password = $1 WHERE EMAIL = $2", newHashedPassword, email)
	if err != nil {
		log.Println("manager.ChangePassword error: " + err.Error())
		return ErrInternal
	}
	return nil
}

// IsEmailRegistered checks if there is an user with a given email
func (manager *Manager) IsEmailRegistered(email string) bool {
	row := manager.Database.QueryRow("SELECT EXISTS(SELECT 1 FROM Users WHERE email = $1)", email)
	var registered bool
	_ = row.Scan(&registered)
	return registered
}

// GetModelFromToken returns an user model based on a JWT token
func (manager *Manager) GetModelFromToken(token string, secret []byte) (*Model, error) {
	claims, err := userauth.GetClaims(token, secret)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return New(claims.Email), nil
}
