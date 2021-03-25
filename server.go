package main

import (
	"encoding/json"
	"github.com/matherique/identity-service/cmd/config"
	"net/http"
)

type DefaultResponse struct {
	Message string
}

func NewServer(port string, config *config.Config) error {
	http.HandleFunc("/", handlerHome)
	return http.ListenAndServe(port, nil)
}

func handlerHome(w http.ResponseWriter, r *http.Request) {
	msg := DefaultResponse{Message: "works!"}

	json.NewEncoder(w).Encode(msg)
}
