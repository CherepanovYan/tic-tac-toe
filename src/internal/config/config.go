package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strings"
)

const (
	CONFIG_FILE = "internal/config/.env"
)

type Config struct {
	Port string
}

func Load() (Config, error) {

	err := godotenv.Load(CONFIG_FILE)
	if err != nil {
		return Config{}, fmt.Errorf("Error loading .env file")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok || strings.TrimSpace(port) == "" {
		return Config{}, fmt.Errorf("PORT is required")
	}

	return Config{Port: strings.TrimSpace(port)}, nil
}
