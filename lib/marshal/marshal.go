package marshal

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
)

type Marshal interface {
	Unmarshal(byte, *interface{}) error
}

type YML struct{}
type JSON struct{}

func (yml *YML) Unmarshal(data []byte, output *interface{}) error {
	if err := yaml.Unmarshal(data, output); err != nil {
		return fmt.Errorf("Error %v", err)
	}
	return nil
}

func (j *JSON) Unmarshal(data []byte, output *interface{}) error {
	if err := json.Unmarshal(data, output); err != nil {
		return fmt.Errorf("Error %v", err)
	}
	return nil
}
