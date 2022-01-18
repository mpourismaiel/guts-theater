package main

import (
	"os"

	"go.uber.org/zap"
	"mpourismaiel.dev/guts/api"
)

func main() {
	logger := zap.NewExample()
	logger.Info("Starting project...")
	defer logger.Sync()

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	address := os.Getenv("ADDRESS")
	if address == "" {
		address = "0.0.0.0"
	}

	dbUser := os.Getenv("DB_USER")
	if address == "" {
		address = "admin"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "password"
	}

	if err := api.New(address, port, dbUser, dbPassword, logger); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
