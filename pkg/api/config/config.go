package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func GetDBConnection() (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", errors.New("Failed to read configuration file")
	}
	return os.Getenv("DB_URI"), nil
}
