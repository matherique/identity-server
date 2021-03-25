package token

import (
	"strings"
	"testing"
	"time"
)

var secret []byte = []byte("1234")
var payload string = "service_test"

func TestCreateToken(t *testing.T) {
	tokenPrefix := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	exp := time.Now().Add(time.Hour * 2).Unix()

	token, err := GenerateToken(payload, secret, exp)

	if err != nil {
		t.Errorf("error in token creation: %v", err)
	}

	if !strings.HasPrefix(token, tokenPrefix) {
		t.Errorf("expect token that start with %s, got %s", token, token)
	}
}

func TestValidateToken(t *testing.T) {
	exp := time.Now().Add(time.Hour * 1).Unix()

	token, err := GenerateToken(payload, secret, exp)

	isValid, err := ValidateToken(token, secret)

	if err != nil {
		t.Errorf("error in token validation: %v", err)
	}

	if !isValid {
		t.Errorf("expect validation return %t, got %t", true, false)
	}
}

func TestGetData(t *testing.T) {
	exp := time.Now().Add(time.Hour * 1).Unix()

	token, err := GenerateToken(payload, secret, exp)

	data, err := GetTokenData(token, secret)

	if err != nil {
		t.Errorf("could not get data from token: %v", err)
	}

	if data != payload {
		t.Errorf("expect data to be %v, got %v", payload, data)
	}

}
