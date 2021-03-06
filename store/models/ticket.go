package models

import (
	"context"
	"fmt"

	"github.com/go-kivik/kivik/v3"
	"github.com/google/uuid"
	"github.com/mpourismaiel/guts-theater/prometheus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// type definition for document
type Ticket struct {
	ID      string   `json:"_id"`
	Rev     string   `json:"_rev,omitempty"`
	GroupId string   `json:"groupId"`
	Seats   []string `json:"seats"`
}

// create views for the model
func (m *Models) ticketCreateModel() error {
	_, err := m.db.Put(context.TODO(), "_design/ticket", map[string]interface{}{
		"id": "_design/ticket",
		"views": map[string]interface{}{
			"ticket-list-all": map[string]string{
				"map": "function (doc) {\n if (!doc._id.match(/^ticket:[^:]+$/)) {\n return;\n }\n emit(doc._id, 1);\n }",
			},
			"ticket-get-by-groupid": map[string]string{
				"map": "function (doc) {\n if (!doc._id.match(/^ticket:[^:]+$/)) {\n return;\n }\n emit(doc.groupId, 1);\n }",
			},
		},
	})
	prometheus.DbCall.WithLabelValues("ticket", "migration").Inc()

	return err
}

// generate id with prefix to indicate type and possible relations
func ticketCreateId() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("ticket:%s", id.String()), nil
}

func (m *Models) TicketSave(t *Ticket) error {
	if t.Rev != "" {
		return fmt.Errorf("failed to save new row due to rev being present: %s", t.Rev)
	}

	id, err := ticketCreateId()
	if err != nil {
		return err
	}
	prometheus.DbCall.WithLabelValues("ticket", "save").Inc()

	t.ID = id
	rev, err := m.db.Put(context.TODO(), t.ID, &t)
	if err != nil {
		return err
	}

	fields := []zapcore.Field{
		zap.String("ticketId", t.ID),
		zap.String("rev", rev),
	}
	m.logger.Debug("Successfully stored ticket", fields...)
	t.Rev = rev
	return nil
}

func (m *Models) TicketDelete(t *Ticket) error {
	if t.Rev == "" {
		return fmt.Errorf("failed to delete ticket (%s) because no rev was provided", t.GroupId)
	}

	rev, err := m.db.Delete(context.TODO(), t.ID, t.Rev)
	if err != nil {
		return err
	}
	prometheus.DbCall.WithLabelValues("ticket", "delete").Inc()

	fields := []zapcore.Field{
		zap.String("ticketId", t.ID),
		zap.String("rev", rev),
	}
	m.logger.Debug("Successfully deleted ticket", fields...)
	t.Rev = rev
	return nil
}

func (m *Models) TicketGetAll() ([]*Ticket, error) {
	docs, err := m.db.Query(context.TODO(), "_design/ticket", "_view/ticket-list-all", kivik.Options{
		"include_docs": true,
	})
	if err != nil {
		return nil, err
	}
	prometheus.DbCall.WithLabelValues("ticket", "query").Inc()

	var result []*Ticket
	for docs.Next() {
		var doc Ticket
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

func (m *Models) TicketGetByGroupId(groupId string) (*Ticket, error) {
	docs, err := m.db.Query(context.TODO(), "_design/ticket", "_view/ticket-get-by-groupid", kivik.Options{
		"include_docs": true,
		"key":          groupId,
	})
	if err != nil {
		return nil, err
	}
	prometheus.DbCall.WithLabelValues("ticket", "query").Inc()

	var doc Ticket
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
