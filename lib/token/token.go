package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	secret []byte
	exp    int64
}

func (t *Token) GenerateToken(payload string) (string, error) {
	claims := jwt.MapClaims{}
	claims["id"] = payload
	claims["exp"] = t.exp

	handler := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := handler.SignedString(t.secret)

	if err != nil {
		return "", fmt.Errorf("can not create token: %v", err)
	}

	return token, nil
}

func (t *Token) verify(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(tk *jwt.Token) (interface{}, error) {
		return []byte(t.secret), nil
	})

}

func (t *Token) Validate(token string) (bool, error) {
	parsed, err := t.verify(token)

	if err != nil {
		return false, fmt.Errorf("invalid token: %v", err)
	}

	return parsed.Valid, nil
}
