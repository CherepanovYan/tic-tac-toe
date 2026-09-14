package config

import (
	"github.com/joho/godotenv"
	"os"
)

const (
	CONFIG_FILE = ".env"
)

type Config struct {
	Port             string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	JwtSecret        string
}

func Load() (Config, error) {

	err := godotenv.Load(CONFIG_FILE)
	if err != nil {
		return Config{}, ErrLoadEnv
	}

	return Config{
		Port:             os.Getenv("APP_PORT"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresPort:     os.Getenv("POSTGRES_PORT"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		JwtSecret:        os.Getenv("JWT_SECRET"),
	}, nil
}
