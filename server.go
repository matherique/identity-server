package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/matherique/identity-service/cmd/config"
	"github.com/matherique/identity-service/lib/token"
)

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

func (s *Server) NewServer() error {
	routes := Routes{config: s.config}

	http.HandleFunc("/", routes.Home)
	http.HandleFunc("/auth", routes.Auth)
	http.HandleFunc("/validate", routes.Validate)

	return http.ListenAndServe(s.port, nil)
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Routes struct {
	config *config.Config
}

func (route *Routes) Home(w http.ResponseWriter, r *http.Request) {
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
}

func (route *Routes) Auth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "not found")

		return
	}

	w.Header().Add("Content-Type", "application/json")
	var resp interface{}

	var data ServerRequest

	json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()

	service, err := route.config.GetService(data.Id)

	if err != nil {
		resp = Response{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("invalid id: %v", err),
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	DAY := time.Hour * 24

	expAccessToken := time.Now().Add(DAY).Unix()
	expRefreshToken := time.Now().Add(DAY * 7).Unix()

	SECRET := []byte("AAAA")

	accessToken, err := token.GenerateToken(service.Id, SECRET, expAccessToken)
	refreshToken, err := token.GenerateToken(service.Id, SECRET, expRefreshToken)

	if err != nil {
		resp = Response{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("error when try to create token: %v", err),
		}

		json.NewEncoder(w).Encode(resp)
		return
	}

	resp = TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	json.NewEncoder(w).Encode(resp)
}

func (route *Routes) Validate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "not found")

		return
	}

	w.Header().Add("Content-Type", "application/json")

	resp := Response{
		Status:  http.StatusOK,
		Message: "not implemented",
	}

	json.NewEncoder(w).Encode(resp)
}
