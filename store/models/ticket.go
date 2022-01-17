package models

import (
	"context"
	"fmt"
	"log"

	"github.com/go-kivik/kivik/v3"
	"github.com/google/uuid"
)

type Ticket struct {
	ID      string   `json:"_id"`
	Rev     string   `json:"_rev,omitempty"`
	GroupId string   `json:"groupId"`
	Seats   []string `json:"seats"`
}

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

	return err
}

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

	t.ID = id
	rev, err := m.db.Put(context.TODO(), t.ID, &t)
	if err != nil {
		return err
	}

	log.Printf("Successfully stored new ticket: %s with revision ID: %s", t.ID, rev)
	t.Rev = rev
	return nil
}

func (m *Models) TicketDelete(t *Ticket) error {
	if t.Rev == "" {
		return fmt.Errorf("failed to delete ticket (%s) because no rev was provided", t.GroupId)
	}

	rev, err := m.db.Delete(context.TODO(), t.ID, t.Rev)
	if err != nil {
		panic(err)
	}

	log.Printf("Successfully deleted ticket: %s. New revision id is: %s", t.GroupId, rev)
	t.Rev = rev
	return nil
}

func (m *Models) TicketGetByGroupId(groupId string) (*Ticket, error) {
	docs, err := m.db.Query(context.TODO(), "_design/ticket", "_view/ticket-get-by-groupid", kivik.Options{
		"include_docs": true,
		"key":          groupId,
	})
	if err != nil {
		return nil, err
	}

	var doc Ticket
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
