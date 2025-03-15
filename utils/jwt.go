package utils

import (
	"chatting-system-backend/objectTypes"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserDetails *objectTypes.LoginCredentials
	jwt.StandardClaims
}

func GenerateJWT(details *objectTypes.LoginCredentials) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserDetails: details,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	return token.SignedString(jwtKey)
}

func ValidateJWT(token string) (bool, error) {
	if token == "" {
		return false, fmt.Errorf("token is empty")
	}
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return false, err
	}
	return tkn.Valid, nil
}

func GetDetailFromJWT(token string) (interface{}, error) {
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	return claims.UserDetails, nil
}
