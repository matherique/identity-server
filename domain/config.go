package config

import (
	marshal "github.com/matherique/identity-service/lib"
)

type Service struct {
	name       string
	depends_on []string
}

type Config struct {
	config   interface{}
	services []Service
}

// SetFromBytes set the config using bytes
func (c *Config) SetFromBytes(data []byte) error {
	var parse marshal.YML

	if err := parse.Unmarshal(data, &rawConfig); err != nil {
		panic(err)
	}
	return nil
}

// SetFromFile set the config using filename
func (c *Config) SetFromFile(filename string) error {

	return nil
}
