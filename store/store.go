package store

import (
	"context"
	"fmt"

	_ "github.com/go-kivik/couchdb/v3"
	kivik "github.com/go-kivik/kivik/v3"
	"github.com/mpourismaiel/guts-theater/config"
	"github.com/mpourismaiel/guts-theater/store/models"
	"go.uber.org/zap"
)

type Orm struct {
	DB     *kivik.DB
	Models *models.Models
	logger *zap.Logger
}

// connects to database and provides access to models package
func New(conf *config.Config, logger *zap.Logger) (*Orm, error) {
	logger.Info("Connect to database")
	client, err := kivik.New("couch", fmt.Sprintf("http://%s:%s@%s:5984/", conf.DbUser, conf.DbPassword, conf.DbHost))
	if err != nil {
		return nil, err
	}

	client.CreateDB(context.TODO(), conf.DbName)
	db := client.DB(context.TODO(), conf.DbName)
	logger.Info("Database connection established")

	m := models.New(db, logger)

	orm := Orm{
		DB:     db,
		Models: m,
		logger: logger,
	}

	return &orm, nil
}
