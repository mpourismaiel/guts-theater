package main

import (
	"os"

	"github.com/mpourismaiel/guts-theater/api"
	"github.com/mpourismaiel/guts-theater/config"
	"github.com/mpourismaiel/guts-theater/prometheus"
	"go.uber.org/zap"
)

func main() {
	logger := zap.NewExample()
	logger.Info("Starting project...")
	defer logger.Sync()

	prometheus.Setup()

	conf := config.Setup()

	if _, err := api.New(conf, logger); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
