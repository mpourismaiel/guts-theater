package models

import (
	"context"
	"fmt"
	"log"

	"github.com/go-kivik/kivik/v3"
	"github.com/google/uuid"
)

type Group struct {
	ID      string `json:"_id"`
	Rev     string `json:"_rev,omitempty"`
	Aisle   bool   `json:"aisle"`
	Rank    string `json:"rank"`
	Count   int    `json:"count"`
	Section string `json:"section"`
}

func (m *Models) groupCreateModel() error {
	_, err := m.db.Put(context.TODO(), "_design/group", map[string]interface{}{
		"id": "_design/group",
		"views": map[string]interface{}{
			"group-list-by-section": map[string]string{
				"map": "function (doc) {\n if (!doc._id.match(/group:[^:]+$/)) {\n return;\n }\n emit(doc.section, 1);\n }",
			},
			"group-get-by-id": map[string]string{
				"map": "function (doc) {\n if (!doc._id.match(/group:[^:]+$/)) {\n return;\n }\n const ids = doc._id.match(/group:([^:]+)$/);\n emit(ids[1], 1);\n }",
			},
		},
	})

	return err
}

func groupCreateId() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("group:%s", id.String()), nil
}

func (m *Models) GroupSave(g *Group) error {
	if g.Rev != "" {
		return fmt.Errorf("failed to save new row due to rev being present: %s", g.Rev)
	}

	id, err := groupCreateId()
	if err != nil {
		return err
	}

	g.ID = id
	rev, err := m.db.Put(context.TODO(), g.ID, &g)
	if err != nil {
		return err
	}

	log.Printf("Successfully stored new group: %s with revision ID: %s", g.ID, rev)
	g.Rev = rev
	return nil
}

func (m *Models) GroupGetAll() ([]*Group, error) {
	docs, err := m.db.Query(context.TODO(), "_design/group", "_view/group-get-by-id", kivik.Options{
		"include_docs": true,
	})
	if err != nil {
		return nil, err
	}

	var result []*Group
	for docs.Next() {
		var doc Group
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

func (m *Models) GroupGetBySection(sectionName string) ([]*Group, error) {
	docs, err := m.db.Query(context.TODO(), "_design/group", "_view/group-list-by-section", kivik.Options{
		"include_docs": true,
		"key":          sectionName,
	})
	if err != nil {
		return nil, err
	}

	var result []*Group
	for docs.Next() {
		var doc Group
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

func (m *Models) GroupGetById(groupId string) (*Group, error) {
	docs, err := m.db.Query(context.TODO(), "_design/group", "_view/group-get-by-id", kivik.Options{
		"include_docs": true,
		"key":          groupId,
	})
	if err != nil {
		return nil, err
	}

	var doc Group
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
