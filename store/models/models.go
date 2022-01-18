package models

import (
	kivik "github.com/go-kivik/kivik/v3"
	"go.uber.org/zap"
)

type Models struct {
	db     *kivik.DB
	logger *zap.Logger
}

// create a models object which contains all types and accessor/modifier methods
// for different documents, also registers prometheus variables which are later
// called when accessing db
func New(db *kivik.DB, logger *zap.Logger) *Models {
	models := Models{
		db:     db,
		logger: logger,
	}

	// create views and any other possible migrations
	models.sectionCreateModel()
	models.rowCreateModel()
	models.seatCreateModel()
	models.groupCreateModel()
	models.ticketCreateModel()

	return &models
}
