package port

import (
	"os"

	"github.com/go-playground/validator"
	yaml "github.com/goccy/go-yaml"
)

// FIXME temp
const DefaultConfigPath = "./ops-agent-ports.yaml"

type Config struct {
	ReservedPorts map[string]uint16 `yaml:"reserved_ports" validate:"required"`
}

func ReadConfig(configPath string) (*Config, error) {
	yamlBytes, err := os.ReadFile(configPath)
	if err != nil {
		config := &Config{ReservedPorts: map[string]uint16{}}
		err = WriteConfig(configPath, config)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
	var config Config
	err = yaml.UnmarshalWithOptions(yamlBytes, &config, yaml.Validator(validator.New()))
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func WriteConfig(configPath string, config *Config) error {
	yamlBytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile(configPath, yamlBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
