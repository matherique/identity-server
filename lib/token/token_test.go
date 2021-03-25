package token

import (
	"strings"
	"testing"
	"time"
)

var secret []byte = []byte("1234")
var idPayload string = "service_test"

func TestCreateToken(t *testing.T) {
	tokenPrefix := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"

	handler := Token{
		secret: secret,
		exp:    time.Now().Add(time.Hour * 2).Unix(),
	}

	token, err := handler.GenerateToken(idPayload)

	if err != nil {
		t.Errorf("error in token creation: %v", err)
	}

	if !strings.HasPrefix(token, tokenPrefix) {
		t.Errorf("expect token that start with %s, got %s", token, token)
	}
}

func TestValidateToken(t *testing.T) {
	handler := Token{
		secret: secret,
		exp:    time.Now().Add(time.Hour * 1).Unix(),
	}

	token, err := handler.GenerateToken(idPayload)

	isValid, err := handler.Validate(token)

	if err != nil {
		t.Errorf("error in token validation: %v", err)
	}

	if !isValid {
		t.Errorf("expect validation return %t, got %t", true, false)
	}
}
