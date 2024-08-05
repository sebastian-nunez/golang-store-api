package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const SEVEN_DAYS_IN_SECONDS int64 = 3600 * 24 * 7

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	JWTExpirationInSeconds int64
	JWTSecret              string
	// When adding new fields, make sure to update `.env.template`
}

// Envs is the global configuration for the application.
var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "8080"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", "1234"),
		DBAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.1"), getEnv("DB_PORT", "3306")),
		DBName:                 getEnv("DB_NAME", "ecommerceDb"),
		JWTExpirationInSeconds: getEnvInt("JWT_EXPIRATION_IN_SECONDS", SEVEN_DAYS_IN_SECONDS),
		JWTSecret:              getEnv("JWT_SECRET", "super-secret"),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
