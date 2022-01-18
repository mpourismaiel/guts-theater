package models

import (
	"context"
	"fmt"

	kivik "github.com/go-kivik/kivik/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Section struct {
	ID        string `json:"_id"`
	Rev       string `json:"_rev,omitempty"`
	Name      string `json:"name"`
	Elevation int    `json:"elevation"`
	Curved    bool   `json:"curved"`
}

func (m *Models) sectionCreateModel() error {
	_, err := m.db.Put(context.TODO(), "_design/section", map[string]interface{}{
		"id": "_design/section",
		"views": map[string]interface{}{
			"section-list-by-name": map[string]string{
				"map": "function (doc) {\n if (doc._id.match(/^section:[^:]+$/)) {\n emit(doc.name, 1);\n }\n }",
			},
		},
	})

	return err
}

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

	fields := []zapcore.Field{
		zap.String("sectionName", s.Name),
		zap.String("rev", rev),
	}
	m.logger.Info("Successfully stored section", fields...)
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

	fields := []zapcore.Field{
		zap.String("sectionName", s.Name),
		zap.String("rev", rev),
	}
	m.logger.Info("Successfully updated section", fields...)
	s.Rev = rev
	return nil
}

func (m *Models) SectionDelete(s *Section) error {
	if s.Rev == "" {
		return fmt.Errorf("failed to delete section (%s) because no rev was provided", s.Name)
	}

	rev, err := m.db.Delete(context.TODO(), s.ID, s.Rev)
	if err != nil {
		panic(err)
	}

	fields := []zapcore.Field{
		zap.String("sectionName", s.Name),
		zap.String("rev", rev),
	}
	m.logger.Info("Successfully deleted section", fields...)
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
