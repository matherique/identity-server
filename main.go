package main

import (
	"fmt"
	"log"
	"os"

	config "github.com/matherique/identity-server/cmd/config"
)

const CONFIG_FILE = "./default.yml"

func main() {
	data, err := os.ReadFile(CONFIG_FILE)

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

	port = fmt.Sprintf(":%s", port)
	server := NewServer(port, &config)

	if err = server.StartServer(); err != nil {
		log.Fatal(err)
	}
}
