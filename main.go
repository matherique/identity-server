package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const CONFIG_FILE = "./default.yml"

type Service struct {
	name       string
	depends_on []string
}

func main() {
	data, err := ioutil.ReadFile(CONFIG_FILE)

	if err != nil {
		panic(err)
	}

	var rawConfig interface{}

	if err := yaml.Unmarshal(data, &rawConfig); err != nil {
		panic(err)
	}

	noTyped, ok := rawConfig.(map[interface{}]interface{})

	if !ok {
		fmt.Errorf("config is not a map")
	}

	var serviceList []Service

	var config map[interface{}]interface{}
	var service Service

	for _, v := range noTyped {
		config, ok = v.(map[interface{}]interface{})
		if !ok {
			continue
		}

		serviceDependsOn := config["depends_on"]

		if serviceDependsOn == nil {
			continue
		}

		var depends_on []string
		listDependsOn := serviceDependsOn.([]interface{})
		for i := 0; i < len(listDependsOn); i++ {
			depends_on = append(depends_on, listDependsOn[i].(string))
		}

		service = Service{
			name:       config["name"].(string),
			depends_on: depends_on,
		}

		serviceList = append(serviceList, service)
	}

	fmt.Println(serviceList)
}
