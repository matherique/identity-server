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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		var resp interface{}
		home := Routes{config: s.config}

		switch r.Method {
		case "GET":
			resp = home.Home(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "not found")

			return
		}

		json.NewEncoder(w).Encode(resp)
	})

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		var resp interface{}

		switch r.Method {
		case "POST":
			resp = routes.Auth(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "not found")

			return
		}

		json.NewEncoder(w).Encode(resp)
	})

	return http.ListenAndServe(s.port, nil)
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type Routes struct {
	config *config.Config
}

func (h *Routes) Home(w http.ResponseWriter, r *http.Request) interface{} {
	return Response{
		Status:  http.StatusOK,
		Message: "works",
	}
}

func (h *Routes) Auth(w http.ResponseWriter, r *http.Request) interface{} {
	var resp interface{}
	var data ServerRequest

	json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()

	service, err := h.config.GetService(data.Id)

	if err != nil {
		resp = Response{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("invalid id: %v", err),
		}
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
	}

	resp = TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return resp
}
