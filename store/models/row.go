package models

import (
	"context"
	"fmt"
	"log"

	kivik "github.com/go-kivik/kivik/v3"
)

type Row struct {
	ID      string `json:"_id"`
	Rev     string `json:"rev,omitempty"`
	Name    string `json:"name"`
	Section string `json:"section"`
}

func (m *Models) rowCreateModel() error {
	_, err := m.db.Put(context.TODO(), "_design/row", map[string]interface{}{
		"id": "_design/row",
		"views": map[string]interface{}{
			"row-list-by-section": map[string]string{
				"map": "function (doc) {\n if (!doc._id.match(/section:[^:]+:row:[^:]+$/)) {\n return;\n }\n emit(doc.section, 1);\n }",
			},
		},
	})

	return err
}

func (m *Models) RowSave(r *Row) error {
	if r.Rev != "" {
		return fmt.Errorf("failed to save new row due to rev being present: %s", r.Rev)
	}

	r.ID = fmt.Sprintf("section:%s:row:%s", r.Section, r.Name)
	rev, err := m.db.Put(context.TODO(), r.ID, &r)
	if err != nil {
		return err
	}

	log.Printf("Successfully stored new row: %s with revision ID: %s", r.Name, rev)
	r.Rev = rev
	return nil
}

func (m *Models) RowUpdate(r *Row) error {
	if r.Rev == "" {
		return fmt.Errorf("failed to update row (%s) because no rev was provided", r.Name)
	}

	rev, err := m.db.Put(context.TODO(), r.ID, &r)
	if err != nil {
		return err
	}

	log.Printf("Successfully updated row: %s with revision ID: %s", r.Name, rev)
	r.Rev = rev
	return nil
}

func (m *Models) RowDelete(r *Row) error {
	if r.Rev == "" {
		return fmt.Errorf("failed to delete row (%s) because no rev was provided", r.Name)
	}

	rev, err := m.db.Delete(context.TODO(), r.ID, r.Rev)
	if err != nil {
		panic(err)
	}

	log.Printf("Successfully deleted row: %s. New revision id is: %s", r.Name, rev)
	r.Rev = rev
	return nil
}

func (m *Models) RowGetBySection(sectionName string) ([]*Row, error) {
	docs, err := m.db.Query(context.TODO(), "_design/row", "_view/row-list-by-section", kivik.Options{
		"include_docs": true,
		"key":          sectionName,
	})
	if err != nil {
		return nil, err
	}

	var result []*Row
	for docs.Next() {
		var doc Row
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
