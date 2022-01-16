package store

import (
	"context"
	"log"

	_ "github.com/go-kivik/couchdb/v3"
	kivik "github.com/go-kivik/kivik/v3"
	"mpourismaiel.dev/guts/store/models"
)

type Orm struct {
	DB     *kivik.DB
	Models *models.Models
}

func New(dbName string) *Orm {
	if dbName == "" {
		dbName = "guts"
	}

	log.Println("Connect to database")
	client, err := kivik.New("couch", "http://admin:password@localhost:5984/")
	if err != nil {
		log.Fatal(err)
	}

	client.CreateDB(context.TODO(), dbName)
	db := client.DB(context.TODO(), dbName)
	log.Println("Database connection established")

	_, err = db.Query(context.TODO(), "_design/result", "_view/result", kivik.Options{
		"include_docs": true,
	})
	if err != nil {
		log.Println("Seeding data")
		// TODO: seed data
	}

	m := models.New(db)

	orm := Orm{
		DB:     db,
		Models: m,
	}

	return &orm
}
