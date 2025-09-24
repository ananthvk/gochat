package config

import (
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Host                      string        `env:"GOCHAT_HOST" envDefault:"127.0.0.1"`
	Port                      int           `env:"GOCHAT_PORT,notEmpty" envDefault:"8000"`
	Env                       string        `env:"GOCHAT_ENV" envDefault:"production"`
	EnableDetailedHealthCheck bool          `env:"GOCHAT_ENABLE_DETAILED_HEALTHCHECK" envDefault:"true"`
	DbDSN                     string        `env:"GOCHAT_DB_DSN,notEmpty"`
	DbPingTimeout             time.Duration `env:"GOCHAT_DB_PING_TIMEOUT,notEmpty" envDefault:"5s"`
}

func LoadEnv() {
	godotenv.Load()
}

func ParseConfig() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return &cfg, err
}
