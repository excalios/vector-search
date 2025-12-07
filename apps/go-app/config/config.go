package config

import (
	"context"
	"go-app/internal/logging"

	"github.com/joho/godotenv"
)

type contextKey string

const (
	AppEnvKey contextKey = "APP_ENVIRONMENT"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		logging.LogInfo(context.Background(), "No env file found")
	}

}
