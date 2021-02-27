package main

import (
	"github.com/BurntSushi/toml"
	"os"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type loggerConfig struct {
	Level string
	// TODO
}

type storageConfig struct {
	Type string
}

type Config struct {
	Logger  loggerConfig
	Storage storageConfig
	// TODO
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

// TODO
