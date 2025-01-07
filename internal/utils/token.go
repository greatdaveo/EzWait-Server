package utils

import (
	"errors"
	"ezwait/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey = []byte("oi3hugj-0987ewh")

// To generate a JWT for a user
func GenerateToken(user *models.User) (string, error) {

	claims := jwt.MapClaims{
		"user": user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}

	// To create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// To sign the token
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// To validates the token and return the claims
func VerifyToken(tokenStr string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	// To extract claims
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, errors.New("invalid claims")
	}

	return &claims, nil
}
