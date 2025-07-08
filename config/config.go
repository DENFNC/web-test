package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPass     string
	DBName     string
	RedisAddr  string
	AdminToken string
}

func Load() *Config {
	_ = godotenv.Load()
	return &Config{
		DBHost:     getenv("DB_HOST", "localhost"),
		DBPort:     getenv("DB_PORT", "5432"),
		DBUser:     getenv("DB_USER", "webuser"),
		DBPass:     getenv("DB_PASSWORD", "webpass"),
		DBName:     getenv("DB_NAME", "websrv"),
		RedisAddr:  getenv("REDIS_ADDR", "localhost:6379"),
		AdminToken: getenv("ADMIN_TOKEN", "superadmin123"),
	}
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
