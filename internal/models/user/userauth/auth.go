package userauth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	tokenIssuer       = "adjsky"
	tokenSubject      = "auth"
	authTokenLifespan = time.Hour * 24
)

// Claims holds user information passed by Authorization HTTP header
type Claims struct {
	Email string
	jwt.StandardClaims
}

// GenerateClaims generates a new JWT token claims
func GenerateClaims(email string) *Claims {
	return &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			Issuer:    tokenIssuer,
			Subject:   tokenSubject,
			ExpiresAt: time.Now().Add(authTokenLifespan).Unix(),
		},
	}
}

// GenerateToken returns a JWT string that is passed to a client
func GenerateToken(claims *Claims, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(secretKey)
	return ss, err
}

// GetClaims decodes a JWT string passed by a client and returns data associated with it if the token is valid
func GetClaims(tokenString string, secretKey []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok {
		if claims.Issuer == tokenIssuer && claims.Subject == tokenSubject {
			return claims, nil
		}
	}
	return nil, errors.New("token has invalid claims")
}
