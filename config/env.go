package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			Log.Info("No .env file found")
		}
	}
}