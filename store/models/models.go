package models

import (
	kivik "github.com/go-kivik/kivik/v3"
)

type Models struct {
	db *kivik.DB
}

func New(db *kivik.DB) *Models {
	models := Models{
		db: db,
	}

	models.sectionCreateModel()
	models.rowCreateModel()
	models.seatCreateModel()

	return &models
}
