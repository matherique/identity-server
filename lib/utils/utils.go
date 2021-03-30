package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/matherique/identity-service/lib/token"
)

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func Authenticate(w http.ResponseWriter, r *http.Request, secret []byte) (bool, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return false, nil
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	isValid, err := token.ValidateToken(tokenString, secret)

	if err != nil {
		return false, err
	}

	if !isValid {
		fmt.Println("invalid token")
		return false, fmt.Errorf("invalid token")
	}

	return true, nil
}

func GenerateTokens(payload string, secret []byte) (interface{}, error) {
	DAY := time.Hour * 24

	expAccessToken := time.Now().Add(DAY).Unix()
	expRefreshToken := time.Now().Add(DAY * 7).Unix()

	accessToken, err := token.GenerateToken(payload, secret, expAccessToken)
	refreshToken, err := token.GenerateToken(payload, secret, expRefreshToken)

	if err != nil {
		return nil, fmt.Errorf("error in token creation: %v", err)
	}

	return TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
