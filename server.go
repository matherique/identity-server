package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/matherique/identity-service/cmd/config"
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		var resp interface{}
		home := Home{config: s.config}

		switch r.Method {
		case "GET":
			resp = home.GET(w, r)

		case "POST":
			resp = home.POST(w, r)
		default:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "not found")

			return
		}

		json.NewEncoder(w).Encode(resp)
	})

	return http.ListenAndServe(s.port, nil)
}

type Home struct {
	config *config.Config
}

func (h *Home) GET(w http.ResponseWriter, r *http.Request) interface{} {
	return Response{
		Status:  http.StatusOK,
		Message: "works",
	}
}

func (h *Home) POST(w http.ResponseWriter, r *http.Request) interface{} {
	var resp Response
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

	resp = Response{
		Status:  http.StatusOK,
		Message: fmt.Sprintf("service name: %v", service.Name),
	}

	return resp
}
