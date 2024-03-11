package config

import (
	"github.com/caarlos0/env"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Port int `env:"PORT" envDefault:"9000"`
}

func Load() Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		logrus.WithError(err).Fatal("Failed to load configuration.")
	}

	return cfg
}
