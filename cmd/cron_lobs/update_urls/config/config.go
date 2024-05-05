package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		ServiceUrl  string `yaml:"service_url"`
		SefviceApi  string `yaml:"service_api"`
		ServicePort string `yaml:"service_port"`
	}
)

func Get(confFile string) (config *Config, err error) {
	rawData, err := os.ReadFile(confFile)
	if err != nil {
		return nil, err
	}

	config = &Config{}

	if err := config.Unmarshal(rawData); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Unmarshal(data []byte) (err error) {

	if err = yaml.Unmarshal(data, c); err != nil {
		return err
	}

	if c.ServiceUrl == "" {
		return fmt.Errorf("config: service_url is empty")
	}
	if c.SefviceApi == "" {
		return fmt.Errorf("config: service_api is empty")
	}
	if c.ServicePort == "" {
		return fmt.Errorf("config: service_port is empty")
	}

	return nil
}
