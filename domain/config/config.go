package config

import (
	"fmt"
	service "github.com/matherique/identity-service/domain/service"
	marshal "github.com/matherique/identity-service/lib"
	"io/ioutil"
)

type Config struct {
	config   interface{}
	services []service.Service
}

// SetFromBytes set the config using bytes
func (c *Config) SetFromBytes(data []byte) error {
	var parse marshal.YML

	var rawConfig interface{}

	if err := parse.Unmarshal(data, &rawConfig); err != nil {
		return fmt.Errorf("Error in marshal data: %v", err.Error())
	}

	noTyped, ok := rawConfig.(map[interface{}]interface{})

	if !ok {
		return fmt.Errorf("config is not a map")
	}

	serviceList, err := convertToService(noTyped)

	if err != nil {
		return fmt.Errorf("could not convert to Service %v", err)
	}

	fmt.Print(serviceList)
	c.services = serviceList

	return nil
}

// SetFromFile set the config using filename
func (c *Config) SetFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return fmt.Errorf("could not read file %v", err)
	}

	if err = c.SetFromBytes(data); err != nil {
		return fmt.Errorf("could not set config from bytes %v", err)
	}

	return nil
}

// convertToService converts map[interface{}]interface{} to service struct
func convertToService(m map[interface{}]interface{}) ([]service.Service, error) {
	var serviceList []service.Service

	var s service.Service

	for k, v := range m {
		config, ok := v.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("could not cast to map[interface{}]interface{}")
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

		s = service.Service{
			Key:        k.(string),
			Name:       config["name"].(string),
			Depends_on: depends_on,
		}

		serviceList = append(serviceList, s)
	}

	return serviceList, nil
}
