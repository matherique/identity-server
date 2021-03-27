package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/matherique/identity-service/cmd/config"
	"github.com/matherique/identity-service/lib/token"
)

var SECRET []byte = []byte("AAAA")

type HandlerRequest = func(http.ResponseWriter, *http.Request)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Server struct {
	config *config.Config
	port   string
}

type ServerRequest struct {
	Id string `json:"id"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (s *Server) NewServer() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "not found")

			return
		}

		resp := Response{
			Status:  http.StatusOK,
			Message: "works",
		}

		json.NewEncoder(w).Encode(resp)

	})

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		if r.Method != "POST" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "not found")

			return
		}

		isAuth, err := Authenticate(w, r, s.config)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, err)

			return
		}

		var resp interface{}

		if isAuth {
			resp = withAuth(w, r, s.config)
		} else {
			resp = withNoAuth(w, r, s.config)
		}

		json.NewEncoder(w).Encode(resp)
	})

	return http.ListenAndServe(s.port, nil)
}

func withNoAuth(w http.ResponseWriter, r *http.Request, config *config.Config) interface{} {
	var resp interface{}

	var data ServerRequest

	json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()

	service, err := config.GetService(data.Id)

	resp = generateTokens(service.Id, SECRET)

	if err != nil {
		resp = Response{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("invalid id: %v", err),
		}
		return resp
	}

	return resp

}

func withAuth(w http.ResponseWriter, r *http.Request, config *config.Config) interface{} {
	var resp interface{}
	var data ServerRequest

	authHeader := r.Header.Get("Authorization")
	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	payload, err := token.GetTokenData(accessToken, SECRET)

	if err != nil {
		resp = Response{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("error in token parser: %v", err),
		}
		return resp
	}

	serviceA, err := config.GetService(payload.(string))

	json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()

	serviceB, err := config.GetService(data.Id)

	if err != nil {
		resp = Response{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("invalid service: %v", err),
		}
		return resp
	}

	if serviceA.IsDependent(serviceB.Id) {
		resp = generateTokens(serviceA.Id, serviceB.Secret)

	}

	return resp

}

func Authenticate(w http.ResponseWriter, r *http.Request, config *config.Config) (bool, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return false, nil
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	isValid, err := token.ValidateToken(tokenString, SECRET)

	if err != nil {
		return false, err
	}

	if !isValid {
		return false, fmt.Errorf("invalid token")
	}

	return true, nil
}

func generateTokens(payload string, secret []byte) interface{} {
	DAY := time.Hour * 24

	expAccessToken := time.Now().Add(DAY).Unix()
	expRefreshToken := time.Now().Add(DAY * 7).Unix()

	accessToken, err := token.GenerateToken(payload, secret, expAccessToken)
	refreshToken, err := token.GenerateToken(payload, secret, expRefreshToken)

	if err != nil {
		return Response{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("error in token creation: %v", err),
		}
	}

	return TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
