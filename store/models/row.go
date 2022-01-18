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
type Row struct {
	ID      string `json:"_id"`
	Rev     string `json:"_rev,omitempty"`
	Name    string `json:"name"`
	Section string `json:"section"`
}

// create views for the model
func (m *Models) rowCreateModel() error {
	_, err := m.db.Put(context.TODO(), "_design/row", map[string]interface{}{
		"id": "_design/row",
		"views": map[string]interface{}{
			"row-list-by-section": map[string]string{
				"map": "function (doc) {\n if (!doc._id.match(/section:[^:]+:row:[^:]+$/)) {\n return;\n }\n emit(doc.section, 1);\n }",
			},
			"row-list-by-name": map[string]string{
				"map": "function (doc) {\n if (!doc._id.match(/section:[^:]+:row:[^:]+$/)) {\n return;\n }\n emit([doc.section, doc.name], 1);\n }",
			},
		},
	})
	prometheus.DbCall.WithLabelValues("row", "migration").Inc()

	return err
}

// generate id with prefix to indicate type and possible relations
func rowCreateId(r *Row) string {
	return fmt.Sprintf("section:%s:row:%s", r.Section, r.Name)
}

func (m *Models) RowSave(r *Row) error {
	if r.Rev != "" {
		return fmt.Errorf("failed to save new row due to rev being present: %s", r.Rev)
	}

	r.ID = rowCreateId(r)
	rev, err := m.db.Put(context.TODO(), r.ID, &r)
	if err != nil {
		return err
	}
	prometheus.DbCall.WithLabelValues("row", "save").Inc()

	fields := []zapcore.Field{
		zap.String("rowName", r.Name),
		zap.String("rev", rev),
	}
	m.logger.Debug("Successfully stored new row", fields...)
	r.Rev = rev
	return nil
}

func (m *Models) RowUpdate(r *Row) error {
	if r.Rev == "" {
		return fmt.Errorf("failed to update row (%s) because no rev was provided", r.Name)
	}

	r.ID = rowCreateId(r)
	rev, err := m.db.Put(context.TODO(), r.ID, &r)
	if err != nil {
		return err
	}
	prometheus.DbCall.WithLabelValues("row", "update").Inc()

	fields := []zapcore.Field{
		zap.String("rowName", r.Name),
		zap.String("rev", rev),
	}
	m.logger.Debug("Successfully updated row", fields...)
	r.Rev = rev
	return nil
}

func (m *Models) RowDelete(r *Row) error {
	if r.Rev == "" {
		return fmt.Errorf("failed to delete row (%s) because no rev was provided", r.Name)
	}

	rev, err := m.db.Delete(context.TODO(), r.ID, r.Rev)
	if err != nil {
		return err
	}
	prometheus.DbCall.WithLabelValues("row", "delete").Inc()

	fields := []zapcore.Field{
		zap.String("rowName", r.Name),
		zap.String("rev", rev),
	}
	m.logger.Debug("Successfully deleted row", fields...)
	r.Rev = rev
	return nil
}

func (m *Models) RowGetByName(sectionName string, rowName string) (*Row, error) {
	docs, err := m.db.Query(context.TODO(), "_design/row", "_view/row-list-by-name", kivik.Options{
		"include_docs": true,
		"key":          []string{sectionName, rowName},
	})
	if err != nil {
		return nil, err
	}
	prometheus.DbCall.WithLabelValues("row", "query").Inc()

	var doc Row
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

func (m *Models) RowGetBySection(sectionName string) ([]*Row, error) {
	docs, err := m.db.Query(context.TODO(), "_design/row", "_view/row-list-by-section", kivik.Options{
		"include_docs": true,
		"key":          sectionName,
	})
	if err != nil {
		return nil, err
	}
	prometheus.DbCall.WithLabelValues("row", "query").Inc()

	var result []*Row
	for docs.Next() {
		var doc Row
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
