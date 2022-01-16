package models

import (
	"context"
	"fmt"
	"log"

	kivik "github.com/go-kivik/kivik/v3"
)

type Seat struct {
	ID      string `json:"_id"`
	Rev     string `json:"rev,omitempty"`
	Row     string `json:"row"`
	Section string `json:"section"`
	Name    string `json:"name"`
	Rank    string `json:"rank"`
	Broken  bool   `json:"broken"`
	Aisle   bool   `json:"aisle"`
}

func (m *Models) seatCreateModel() error {
	_, err := m.db.Put(context.TODO(), "_design/seat", map[string]interface{}{
		"id": "_design/seat",
		"views": map[string]interface{}{
			"seat-list-all": map[string]string{
				"map": "function (doc) {\n if (!doc._id.match(/^section:[^:]+:row:[^:]+:seat:[^:]+$/)) {\n return;\n }\n emit([doc.section, doc.row, doc.name], 1);\n }",
			},
			"seat-list-by-row": map[string]string{
				"map": "function (doc) {\n if (!doc._id.match(/^section:[^:]+:row:[^:]+:seat:[^:]+$/)) {\n return;\n }\n emit([doc.section, doc.row], 1);\n }",
			},
			"seat-list-by-section": map[string]string{
				"map": "function (doc) {\n if (!doc._id.match(/^section:[^:]+:row:[^:]+:seat:[^:]+$/)) {\n return;\n }\n emit(doc.section, 1);\n }",
			},
		},
	})

	return err
}

func (m *Models) SeatSave(s *Seat) error {
	if s.Rev != "" {
		return fmt.Errorf("failed to save new seat due to rev being present: %s", s.Rev)
	}

	s.ID = fmt.Sprintf("section:%s:row:%s:seat:%s", s.Section, s.Row, s.Name)
	rev, err := m.db.Put(context.TODO(), s.ID, &s)
	if err != nil {
		return err
	}

	log.Printf("Successfully stored new seat: %s with revision ID: %s", s.Name, rev)
	s.Rev = rev
	return nil
}

func (m *Models) SeatUpdate(s *Seat) error {
	if s.Rev == "" {
		return fmt.Errorf("failed to update seat (%s) because no rev was provided", s.Name)
	}

	rev, err := m.db.Put(context.TODO(), s.ID, &s)
	if err != nil {
		return err
	}

	log.Printf("Successfully updated seat: %s with revision ID: %s", s.Name, rev)
	s.Rev = rev
	return nil
}

func (m *Models) SeatDelete(s *Seat) error {
	if s.Rev == "" {
		return fmt.Errorf("failed to delete seat (%s) because no rev was provided", s.Name)
	}

	rev, err := m.db.Delete(context.TODO(), s.ID, s.Rev)
	if err != nil {
		panic(err)
	}

	log.Printf("Successfully deleted seat: %s. New revision id is: %s", s.Name, rev)
	s.Rev = rev
	return nil
}

func (m *Models) SeatGetByRow(sectionName string, rowName string) ([]*Seat, error) {
	docs, err := m.db.Query(context.TODO(), "_design/seat", "_view/seat-list-by-row", kivik.Options{
		"include_docs": true,
		"key":          []string{sectionName, rowName},
	})
	if err != nil {
		return nil, err
	}

	var result []*Seat
	for docs.Next() {
		var doc Seat
		if err := docs.ScanDoc(&doc); err != nil {
			panic(err)
		}
		result = append(result, &doc)
	}

	if docs.Err() != nil {
		panic(docs.Err())
	}

	return result, nil
}

func (m *Models) SeatGetBySection(sectionName string) ([]*Seat, error) {
	docs, err := m.db.Query(context.TODO(), "_design/seat", "_view/seat-list-by-section", kivik.Options{
		"include_docs": true,
		"key":          sectionName,
	})
	if err != nil {
		return nil, err
	}

	var result []*Seat
	for docs.Next() {
		var doc Seat
		if err := docs.ScanDoc(&doc); err != nil {
			panic(err)
		}
		result = append(result, &doc)
	}

	if docs.Err() != nil {
		panic(docs.Err())
	}

	return result, nil
}

func (m *Models) SeatGetAll() ([]*Seat, error) {
	docs, err := m.db.Query(context.TODO(), "_design/seat", "_view/seat-list-all", kivik.Options{
		"include_docs": true,
	})
	if err != nil {
		return nil, err
	}

	var result []*Seat
	for docs.Next() {
		var doc Seat
		if err := docs.ScanDoc(&doc); err != nil {
			panic(err)
		}
		result = append(result, &doc)
	}

	if docs.Err() != nil {
		panic(docs.Err())
	}

	return result, nil
}
