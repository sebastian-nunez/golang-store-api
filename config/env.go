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
	DbUser                 string
	DbPassword             string
	DbAddress              string
	DbName                 string
	JwtExpirationInSeconds int64
	JwtSecret              string
}

// Envs is the global configuration for the application.
var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "8080"),
		DbUser:                 getEnv("DB_USER", "root"),
		DbPassword:             getEnv("DB_PASSWORD", "1234"),
		DbAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.1"), getEnv("DB_PORT", "3306")),
		DbName:                 getEnv("DB_NAME", "ecommerceDb"),
		JwtExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", SEVEN_DAYS_IN_SECONDS),
		JwtSecret:              getEnv("JWT_SECRET", "super-secret"),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
