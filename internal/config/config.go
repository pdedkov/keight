package config

import (
	"github.com/kelseyhightower/envconfig"
)

const (
	prefix = ""
)

type Config struct {
	API     *API
	Log     *Log
	Storage *Storage
	Redis   *Redis
}

// New process and creates new app config
func New() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process(prefix, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
