package models

import (
	"context"
	"fmt"

	kivik "github.com/go-kivik/kivik/v3"
	"github.com/mpourismaiel/guts-theater/prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// type definition for document
type Seat struct {
	ID      string `json:"_id"`
	Rev     string `json:"_rev,omitempty"`
	Row     string `json:"row"`
	Section string `json:"section"`
	Name    string `json:"name"`
	Rank    string `json:"rank"`
	Broken  bool   `json:"broken"`
	Aisle   bool   `json:"aisle"`
}

// create views for the model
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
	prometheus.DbCall.WithLabelValues("seat", "migration").Inc()

	return err
}

// generate id with prefix to indicate type and possible relations
func seatCreateId(s *Seat) string {
	return fmt.Sprintf("section:%s:row:%s:seat:%s", s.Section, s.Row, s.Name)
}

func (m *Models) SeatSave(s *Seat) error {
	if s.Rev != "" {
		return fmt.Errorf("failed to save new seat due to rev being present: %s", s.Rev)
	}

	s.ID = seatCreateId(s)
	rev, err := m.db.Put(context.TODO(), s.ID, &s)
	if err != nil {
		return err
	}
	prometheus.DbCall.WithLabelValues("seat", "save").Inc()

	fields := []zapcore.Field{
		zap.String("seatName", s.Name),
		zap.String("rev", rev),
	}
	m.logger.Debug("Successfully stored new seat", fields...)
	s.Rev = rev
	return nil
}

func (m *Models) SeatUpdate(s *Seat) error {
	if s.Rev == "" {
		return fmt.Errorf("failed to update seat (%s) because no rev was provided", s.Name)
	}

	s.ID = seatCreateId(s)
	rev, err := m.db.Put(context.TODO(), s.ID, &s)
	if err != nil {
		return err
	}
	prometheus.DbCall.WithLabelValues("seat", "update").Inc()

	fields := []zapcore.Field{
		zap.String("seatName", s.Name),
		zap.String("rev", rev),
	}
	m.logger.Debug("Successfully updated seat", fields...)
	s.Rev = rev
	return nil
}

func (m *Models) SeatDelete(s *Seat) error {
	if s.Rev == "" {
		return fmt.Errorf("failed to delete seat (%s) because no rev was provided", s.Name)
	}

	rev, err := m.db.Delete(context.TODO(), s.ID, s.Rev)
	if err != nil {
		return err
	}
	prometheus.DbCall.WithLabelValues("seat", "delete").Inc()

	fields := []zapcore.Field{
		zap.String("seatName", s.Name),
		zap.String("rev", rev),
	}
	m.logger.Debug("Successfully deleted seat", fields...)
	s.Rev = rev
	return nil
}

func (m *Models) SeatGetByName(sectionName string, rowName string, seatName string) (*Seat, error) {
	docs, err := m.db.Query(context.TODO(), "_design/seat", "_view/seat-list-all", kivik.Options{
		"include_docs": true,
		"key":          []string{sectionName, rowName, seatName},
	})
	if err != nil {
		return nil, err
	}
	prometheus.DbCall.WithLabelValues("seat", "query").Inc()

	var doc Seat
	for docs.Next() {
		if err := docs.ScanDoc(&doc); err != nil {
			return nil, err
		}
	}

	if docs.Err() != nil {
		return nil, err
	}

	return &doc, nil
}

func (m *Models) SeatGetByRow(sectionName string, rowName string) ([]*Seat, error) {
	docs, err := m.db.Query(context.TODO(), "_design/seat", "_view/seat-list-by-row", kivik.Options{
		"include_docs": true,
		"key":          []string{sectionName, rowName},
	})
	if err != nil {
		return nil, err
	}
	prometheus.DbCall.WithLabelValues("seat", "query").Inc()

	var result []*Seat
	for docs.Next() {
		var doc Seat
		if err := docs.ScanDoc(&doc); err != nil {
			return nil, err
		}
		result = append(result, &doc)
	}

	if docs.Err() != nil {
		return nil, err
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
	prometheus.DbCall.WithLabelValues("seat", "query").Inc()

	var result []*Seat
	for docs.Next() {
		var doc Seat
		if err := docs.ScanDoc(&doc); err != nil {
			return nil, err
		}
		result = append(result, &doc)
	}

	if docs.Err() != nil {
		return nil, err
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
	prometheus.DbCall.WithLabelValues("seat", "query").Inc()

	var result []*Seat
	for docs.Next() {
		var doc Seat
		if err := docs.ScanDoc(&doc); err != nil {
			return nil, err
		}
		result = append(result, &doc)
	}

	if docs.Err() != nil {
		return nil, err
	}

	return result, nil
}
