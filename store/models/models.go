package models

import (
	kivik "github.com/go-kivik/kivik/v3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

type Models struct {
	db     *kivik.DB
	logger *zap.Logger
}

var (
	dbCall = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:        "guts_theater_db_call_total",
		Help:        "The total number of calls to database",
		ConstLabels: prometheus.Labels{"service": "guts"},
	}, []string{"model", "action"})
)

func New(db *kivik.DB, logger *zap.Logger) *Models {
	models := Models{
		db:     db,
		logger: logger,
	}
	prometheus.MustRegister(dbCall)

	models.sectionCreateModel()
	models.rowCreateModel()
	models.seatCreateModel()
	models.groupCreateModel()
	models.ticketCreateModel()

	return &models
}
