package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	App    AppSection    `yaml:"app"`
	Server ServerSection `yaml:"server"`
}

type AppSection struct {
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
	Debug       bool   `yaml:"debug"`
	Version     string `yaml:"version"`
	URL         string `yaml:"url"`
}

type ServerSection struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

func Load() (*AppConfig, error) {
	return LoadAppConfig("")
}

func LoadAppConfig(configPath string) (*AppConfig, error) {
	if configPath == "" {
		configPath = "config/app.yaml"
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}
