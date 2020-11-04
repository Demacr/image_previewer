package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	Host          string `env:"HOST" envDefault:"localhost"`
	Port          int    `env:"PORT" envDefault:"8080"`
	CacheCapacity int    `env:"CACHE_CAPACITY" envDefault:"10"`
}

func Configure() (*Config, error) {
	config := &Config{}
	if err := env.Parse(config); err != nil {
		return nil, err
	}
	return config, nil
}
