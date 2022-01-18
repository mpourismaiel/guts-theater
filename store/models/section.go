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
type Section struct {
	ID        string `json:"_id"`
	Rev       string `json:"_rev,omitempty"`
	Name      string `json:"name"`
	Elevation int    `json:"elevation"`
	Curved    bool   `json:"curved"`
}

// create views for the model
func (m *Models) sectionCreateModel() error {
	_, err := m.db.Put(context.TODO(), "_design/section", map[string]interface{}{
		"id": "_design/section",
		"views": map[string]interface{}{
			"section-list-by-name": map[string]string{
				"map": "function (doc) {\n if (doc._id.match(/^section:[^:]+$/)) {\n emit(doc.name, 1);\n }\n }",
			},
		},
	})
	prometheus.DbCall.WithLabelValues("section", "migration").Inc()

	return err
}

// generate id with prefix to indicate type and possible relations
func sectionCreateId(s *Section) string {
	return fmt.Sprintf("section:%s", s.Name)
}

func (m *Models) SectionSave(s *Section) error {
	if s.Rev != "" {
		return fmt.Errorf("failed to save new section due to rev being present: %s", s.Rev)
	}

	s.ID = sectionCreateId(s)
	rev, err := m.db.Put(context.TODO(), s.ID, &s)
	if err != nil {
		return err
	}
	prometheus.DbCall.WithLabelValues("section", "save").Inc()

	fields := []zapcore.Field{
		zap.String("sectionName", s.Name),
		zap.String("rev", rev),
	}
	m.logger.Debug("Successfully stored section", fields...)
	s.Rev = rev
	return nil
}

func (m *Models) SectionUpdate(s *Section) error {
	if s.Rev == "" {
		return fmt.Errorf("failed to update section (%s) because no rev was provided", s.Name)
	}

	s.ID = sectionCreateId(s)
	rev, err := m.db.Put(context.TODO(), s.ID, &s)
	if err != nil {
		return err
	}
	prometheus.DbCall.WithLabelValues("section", "update").Inc()

	fields := []zapcore.Field{
		zap.String("sectionName", s.Name),
		zap.String("rev", rev),
	}
	m.logger.Debug("Successfully updated section", fields...)
	s.Rev = rev
	return nil
}

func (m *Models) SectionDelete(s *Section) error {
	if s.Rev == "" {
		return fmt.Errorf("failed to delete section (%s) because no rev was provided", s.Name)
	}

	rev, err := m.db.Delete(context.TODO(), s.ID, s.Rev)
	if err != nil {
		return err
	}
	prometheus.DbCall.WithLabelValues("section", "delete").Inc()

	fields := []zapcore.Field{
		zap.String("sectionName", s.Name),
		zap.String("rev", rev),
	}
	m.logger.Debug("Successfully deleted section", fields...)
	s.Rev = rev
	return nil
}

func (m *Models) SectionGetAll() ([]*Section, error) {
	docs, err := m.db.Query(context.TODO(), "_design/section", "_view/section-list-by-name", kivik.Options{
		"include_docs": true,
	})
	if err != nil {
		return nil, err
	}
	prometheus.DbCall.WithLabelValues("section", "query").Inc()

	var result []*Section
	for docs.Next() {
		var doc Section
		if err := docs.ScanDoc(&doc); err != nil {
			return nil, err
		}
		result = append(result, &doc)
	}

	if docs.Err() != nil {
		return nil, docs.Err()
	}

	return result, nil
}

func (m *Models) SectionGetByName(name string) (*Section, error) {
	docs, err := m.db.Query(context.TODO(), "_design/section", "_view/section-list-by-name", kivik.Options{
		"include_docs": true,
		"key":          name,
	})
	if err != nil {
		return nil, err
	}
	prometheus.DbCall.WithLabelValues("section", "query").Inc()

	var doc Section
	for docs.Next() {
		if err := docs.ScanDoc(&doc); err != nil {
			return nil, err
		}
	}

	if docs.Err() != nil {
		return nil, docs.Err()
	}

	return &doc, nil
}
