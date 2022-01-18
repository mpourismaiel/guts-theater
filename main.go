package main

import (
	"os"

	"github.com/mpourismaiel/guts-theater/api"
	"go.uber.org/zap"
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
		address = ""
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "admin"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "password"
	}

	if err := api.New(address, port, dbHost, dbUser, dbPassword, logger); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
