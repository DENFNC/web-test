package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig  *DatabaseConfig `env:",init"`
	AppConfig *AppConfig      `env:",init"`
}

type AppConfig struct {
	URL string `env:"APP_URL"`
}

type DatabaseConfig struct {
	URL               string        `env:"DATABASE_URL,required"`
	MaxConns          int32         `env:"DATABASE_MAX_CONNS" envDefault:"25"`
	MinConns          int32         `env:"DATABASE_MIN_CONNS" envDefault:"5"`
	MaxConnLifeTime   time.Duration `env:"DATABASE_MAX_CONN_LIFE_TIME" envDefault:"30m"`
	MaxConnIdleTime   time.Duration `env:"DATABASE_MAX_CONN_IDLE_TIME" envDefault:"5m"`
	HealthCheckPeriod time.Duration `env:"DATABASE_HEALTH_CHECK_PERIOD" envDefault:"1m"`
}

type Cache struct{}

func LoadConfig(log *slog.Logger, path string) *Config {
	const op = "config.LoadConfig"

	log = log.With("op", op)

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			log.Warn(
				"File not found, skipping .env load",
				slog.String("path", path),
			)
		} else {
			log.Error(
				"Error checking file",
				slog.String("err", err.Error()),
			)
		}
	} else {
		if err := godotenv.Load(path); err != nil {
			log.Error(
				"Error reading file",
				slog.String("err", err.Error()),
			)
		}
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Error(
			"Error parsing variables into structure",
			slog.String("err", err.Error()),
		)
		panic(err)
	}

	return &cfg
}
