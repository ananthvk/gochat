package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Host                      string `env:"GOCHAT_HOST" envDefault:"127.0.0.1"`
	Port                      int    `env:"GOCHAT_PORT,notEmpty" envDefault:"8000"`
	Env                       string `env:"GOCHAT_ENV" envDefault:"production"`
	EnableDetailedHealthCheck bool   `env:"GOCHAT_ENABLE_DETAILED_HEALTHCHECK" envDefault:"true"`
}

func LoadEnv() {
	godotenv.Load()
}

func ParseConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return &cfg, err
}
