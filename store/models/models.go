package models

import (
	kivik "github.com/go-kivik/kivik/v3"
	"go.uber.org/zap"
)

type Models struct {
	db     *kivik.DB
	logger *zap.Logger
}

func New(db *kivik.DB, logger *zap.Logger) *Models {
	models := Models{
		db:     db,
		logger: logger,
	}

	models.sectionCreateModel()
	models.rowCreateModel()
	models.seatCreateModel()
	models.groupCreateModel()
	models.ticketCreateModel()

	return &models
}
