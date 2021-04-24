package main

import (
	"os"

	"github.com/BurntSushi/toml"
)

type loggerConfig struct {
	Level string
	Path  string
}

type storageConfig struct {
	Type string
	URI  string
}

type httpConfig struct {
	Host string
	Port string
}

type Config struct {
	Logger  loggerConfig
	Storage storageConfig
	Server  httpConfig
}

func NewConfig() (*Config, error) {
	_, err := os.Stat(configFile)
	if err != nil {
		return nil, err
	}

	cfg := Config{}

	if _, err := toml.DecodeFile(configFile, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
