package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwt_key = os.Getenv("JWT_KEY")

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func GenerateJWTToken(userID uuid.UUID) (*TokenPair, error) {
	AccessClaims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(15 * time.Minute).Unix(),
	}
	AccessString := jwt.NewWithClaims(jwt.SigningMethodHS256, AccessClaims)
	AccesToken, err := AccessString.SignedString(jwt_key)
	if err != nil {
		return nil, err
	}
	RefreshClaims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(3 * 30 * 24 * time.Hour).Unix(),
	}
	RefreshString := jwt.NewWithClaims(jwt.SigningMethodHS256, RefreshClaims)
	RefreshToken, err := RefreshString.SignedString(jwt_key)
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  AccesToken,
		RefreshToken: RefreshToken,
	}, nil
}
