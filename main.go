package main

import (
	"os"
	"time"

	"github.com/getsentry/sentry-go"
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

	if conf.SentryDns != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn: conf.SentryDns,
		}); err != nil {
			sentry.CaptureException(err)
			logger.Fatal(err.Error())
		}
		defer sentry.Flush(2 * time.Second)
	}

	if _, err := api.New(conf, logger); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
