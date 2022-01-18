package models

import (
	"context"
	"fmt"

	"github.com/go-kivik/kivik/v3"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// type definition for document
type Group struct {
	ID      string `json:"_id"`
	Rev     string `json:"_rev,omitempty"`
	Aisle   bool   `json:"aisle"`
	Rank    string `json:"rank"`
	Count   int    `json:"count"`
	Section string `json:"section"`
}

// create views for the model
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
	dbCall.WithLabelValues("group", "migration").Inc()

	return err
}

// generate id with prefix to indicate type and possible relations
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
	dbCall.WithLabelValues("group", "save").Inc()

	g.ID = id
	rev, err := m.db.Put(context.TODO(), g.ID, &g)
	if err != nil {
		return err
	}

	fields := []zapcore.Field{
		zap.String("groupId", g.ID),
		zap.String("rev", rev),
	}
	m.logger.Debug("Successfully stored new group", fields...)
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
	dbCall.WithLabelValues("group", "query").Inc()

	var result []*Group
	for docs.Next() {
		var doc Group
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

func (m *Models) GroupGetBySection(sectionName string) ([]*Group, error) {
	docs, err := m.db.Query(context.TODO(), "_design/group", "_view/group-list-by-section", kivik.Options{
		"include_docs": true,
		"key":          sectionName,
	})
	if err != nil {
		return nil, err
	}
	dbCall.WithLabelValues("group", "query").Inc()

	var result []*Group
	for docs.Next() {
		var doc Group
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

func (m *Models) GroupGetById(groupId string) (*Group, error) {
	docs, err := m.db.Query(context.TODO(), "_design/group", "_view/group-get-by-id", kivik.Options{
		"include_docs": true,
		"key":          groupId,
	})
	if err != nil {
		return nil, err
	}
	dbCall.WithLabelValues("group", "query").Inc()

	var doc Group
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
