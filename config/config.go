package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		App *AppConfig `yaml:"app"`
		DB  *DBConfig  `yaml:"db"`
	}
	AppConfig struct {
		Port string `yaml:"port"`
	}
	DBConfig struct {
		Host       string `yaml:"host"`
		Port       string `yaml:"port"`
		User       string `yaml:"user"`
		Name       string `yaml:"name"`
		SearchPath string `yaml:"search_path"`
		Password   string `yaml:"password"`
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
	c.DB = &DBConfig{}
	c.App = &AppConfig{}

	if err = yaml.Unmarshal(data, c); err != nil {
		return err
	}

	if c.App.Port == "" {
		return fmt.Errorf("config: app.port is empty")
	}
	if c.DB.Host == "" {
		return fmt.Errorf("config: db.host is empty")
	}
	if c.DB.Port == "" {
		return fmt.Errorf("config: db.port is empty")
	}
	if c.DB.User == "" {
		return fmt.Errorf("config: db.user is empty")
	}
	if c.DB.Name == "" {
		return fmt.Errorf("config: db.name is empty")
	}
	if c.DB.SearchPath == "" {
		return fmt.Errorf("config: db.search_path is empty")
	}
	if c.DB.Password == "" {
		return fmt.Errorf("config: db.password is empty")
	}

	return nil
}
