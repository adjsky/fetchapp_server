package user

import "github.com/adjsky/fetchapp_server/internal/models/user/userauth"

// Model is an user data representation
type Model struct {
	Email string
}

// New returns an user model
func New(email string) *Model {
	return &Model{
		Email: email,
	}
}

// GetAuthToken returns a JWT token for authentication
func (model *Model) GetAuthToken(secret []byte) (string, error) {
	claims := userauth.GenerateClaims(model.Email)
	token, err := userauth.GenerateToken(claims, secret)
	if err != nil {
		return "", ErrInternal
	}
	return token, nil
}
