package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	config "github.com/matherique/identity-service/cmd/config"
)

const CONFIG_FILE = "./default.yml"

func main() {
	data, err := ioutil.ReadFile(CONFIG_FILE)

	if err != nil {
		log.Fatal(err)
	}

	secret := os.Getenv("SECRET")

	if secret == "" {
		panic(fmt.Errorf("missing environment variable SECRET"))
	}

	config := config.Config{
		Secret: []byte(secret),
	}

	err = config.SetFromBytes(data)

	if err != nil {
		log.Println(err)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	server := Server{
		port:   fmt.Sprintf(":%s", port),
		config: &config,
	}

	log.Fatal(server.NewServer())
}
