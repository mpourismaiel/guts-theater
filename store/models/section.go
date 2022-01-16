package models

import (
	"context"
	"fmt"
	"log"

	kivik "github.com/go-kivik/kivik/v3"
)

type Section struct {
	ID        string `json:"_id"`
	Rev       string `json:"rev,omitempty"`
	Name      string `json:"name"`
	Elevation int    `json:"elevation"`
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

func (m *Models) SectionSave(s *Section) error {
	if s.Rev != "" {
		return fmt.Errorf("failed to save new section due to rev being present: %s", s.Rev)
	}

	s.ID = fmt.Sprintf("section:%s", s.Name)
	rev, err := m.db.Put(context.TODO(), s.ID, &s)
	if err != nil {
		return err
	}

	log.Printf("Successfully stored new section: %s with revision ID: %s", s.Name, rev)
	s.Rev = rev
	return nil
}

func (m *Models) SectionUpdate(s *Section) error {
	if s.Rev == "" {
		return fmt.Errorf("failed to update section (%s) because no rev was provided", s.Name)
	}

	rev, err := m.db.Put(context.TODO(), s.ID, &s)
	if err != nil {
		return err
	}

	log.Printf("Successfully updated section: %s with revision ID: %s", s.Name, rev)
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

	log.Printf("Successfully deleted section: %s. New revision id is: %s", s.Name, rev)
	s.Rev = rev
	return nil
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
			panic(err)
		}
	}

	if docs.Err() != nil {
		panic(docs.Err())
	}

	return &doc, nil
}
