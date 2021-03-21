package main

import (
	config "github.com/matherique/identity-service/domain/config"
	"io/ioutil"
)

const CONFIG_FILE = "./default.yml"

func main() {
	data, err := ioutil.ReadFile(CONFIG_FILE)

	if err != nil {
		panic(err)
	}

	config := config.Config{}
	err = config.SetFromBytes(data)

	if err != nil {
		panic(err)
	}

}
