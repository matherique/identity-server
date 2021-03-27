package token

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken return a token generated
func GenerateToken(payload string, secret []byte, exp int64) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = payload
	claims["exp"] = exp

	handler := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := handler.SignedString(secret)

	if err != nil {
		return "", fmt.Errorf("can not create token: %v", err)
	}

	return token, nil
}

func verify(token string, secret []byte) (*jwt.Token, error) {
	return jwt.Parse(token, func(tk *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}

// ValidateToken validate a token
func ValidateToken(token string, secret []byte) (bool, error) {
	tokenParsed, err := verify(token, secret)

	if err != nil {
		return false, fmt.Errorf("could not parse token: %v", err)
	}

	return tokenParsed.Valid, nil
}

// GetTokenData parse a token and get data
func GetTokenData(token string, secret []byte) (interface{}, error) {
	tokenParsed, err := verify(token, secret)

	if err != nil {
		return false, err
	}

	payload, ok := tokenParsed.Claims.(jwt.MapClaims)["id"].(string)

	if !ok && !tokenParsed.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return payload, nil
}
