package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/matherique/identity-server/cmd/config"
	"github.com/matherique/identity-server/lib/token"
	"github.com/matherique/identity-server/lib/utils"
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

func NewServer(port string, config *config.Config) *Server {
	return &Server{
		port:   port,
		config: config,
	}
}
func (s *Server) StartServer() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "not found")
			return
		}

		w.WriteHeader(http.StatusOK)
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

		isAuth, err := utils.Authenticate(w, r, s.config.Secret)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			json.NewEncoder(w).Encode(Response{
				Status:  http.StatusUnauthorized,
				Message: fmt.Sprintf("%v", err),
			})
		}

		var resp interface{}
		var data ServerRequest

		json.NewDecoder(r.Body).Decode(&data)
		defer r.Body.Close()

		if isAuth {

			authHeader := r.Header.Get("Authorization")
			accessToken := strings.TrimPrefix(authHeader, "Bearer ")

			payload, err := token.GetTokenData(accessToken, s.config.Secret)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				resp = Response{
					Status:  http.StatusInternalServerError,
					Message: fmt.Sprintf("error in token parser: %v", err),
				}
				json.NewEncoder(w).Encode(resp)
			}

			resp, err = withAuth(payload.(string), data.Id, s.config)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				resp = Response{
					Status:  http.StatusInternalServerError,
					Message: fmt.Sprintf("%v", err),
				}

				json.NewEncoder(w).Encode(resp)
			}

		} else {
			resp, err = withNoAuth(data.Id, s.config)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				resp = Response{
					Status:  http.StatusInternalServerError,
					Message: fmt.Sprintf("%v", err),
				}

				json.NewEncoder(w).Encode(resp)

			}
		}

		json.NewEncoder(w).Encode(resp)
	})

	return http.ListenAndServe(s.port, nil)
}

func withNoAuth(id string, config *config.Config) (interface{}, error) {
	service, err := config.GetService(id)

	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}

	return utils.GenerateTokens(service.Id, config.Secret)
}

func withAuth(id string, target string, config *config.Config) (interface{}, error) {
	serviceA, err := config.GetService(id)

	if err != nil {
		return nil, fmt.Errorf("invalid service: %v", err)
	}

	serviceB, err := config.GetService(target)

	if err != nil {
		return nil, fmt.Errorf("invalid service: %v", err)
	}

	if !serviceA.IsDependent(serviceB.Id) {
		return nil, fmt.Errorf("%s not depends on %s", serviceA.Name, serviceB.Name)
	}

	return serviceA.GenerateToken(serviceB)
}
