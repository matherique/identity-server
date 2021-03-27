package main

import (
	config "github.com/matherique/identity-service/cmd/config"
	"io/ioutil"
	"log"
)

const CONFIG_FILE = "./default.yml"

func main() {
	data, err := ioutil.ReadFile(CONFIG_FILE)

	if err != nil {
		log.Fatal(err)
	}

	config := config.Config{}
	err = config.SetFromBytes(data)

	if err != nil {
		log.Println(err)
	}

	//s, err := config.GetService("service1")

	if err != nil {
		log.Println(err)
	}

	server := Server{
		port:   ":8000",
		config: &config,
	}

	log.Fatal(server.NewServer())
}
