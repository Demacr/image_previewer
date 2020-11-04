package config

import (
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
)

type Config struct {
	Host          string `env:"HOST" envDefault:"localhost"`
	Port          int    `env:"PORT" envDefault:"8080"`
	CacheCapacity int    `env:"CACHE_CAPACITY" envDefault:"10"`
}

func Configure() (*Config, error) {
	config := &Config{}
	if err := env.Parse(config); err != nil {
		return nil, errors.Wrap(err, "error during parsing env variables")
	}
	return config, nil
}
