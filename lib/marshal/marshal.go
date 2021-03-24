package marshal

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

type Marshal interface {
	Unmarshal(byte, *interface{}) error
}

type YML struct{}

func (yml *YML) Unmarshal(data []byte, output *interface{}) error {
	if err := yaml.Unmarshal(data, output); err != nil {
		return fmt.Errorf("Error %v", err)
	}
	return nil
}
