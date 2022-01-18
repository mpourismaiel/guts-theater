package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"mpourismaiel.dev/guts/server"
	"mpourismaiel.dev/guts/store"
)

type ApiServer struct {
	server *server.Server
	store  *store.Orm
	logger *zap.Logger
}

func New(address string, port string, dbUser string, dbPassword string, logger *zap.Logger) error {
	s, err := server.New(address, port, logger)
	if err != nil {
		return fmt.Errorf("failed to create server for API: %v", err)
	}

	o, err := store.New("guts", dbUser, dbPassword, logger)
	if err != nil {
		return err
	}

	api := ApiServer{
		server: s,
		store:  o,
		logger: logger,
	}

	api.server.Router().Route("/", func(r chi.Router) {
		r.Get("/seats", api.fetchSeats())
		r.Get("/seats/{section}", api.fetchSeatsBySection())

		r.Get("/section", api.fetchSections())
		r.Post("/section", api.createSection())
		r.Put("/section/{section}", api.updateSection())
		r.Delete("/section/{section}", api.deleteSection())

		r.Get("/section/{section}/rows", api.fetchRowsBySection())
		r.Post("/section/{section}/row", api.createRow())
		r.Delete("/section/{section}/row/{row}", api.deleteRow())

		r.Get("/section/{section}/seats", api.fetchSeatsBySection())
		r.Post("/section/{section}/row/{row}/seat", api.createSeats())
		r.Put("/section/{section}/row/{row}/seat/{seat}", api.updateSeat())
		r.Delete("/section/{section}/row/{row}/seat/{seat}", api.deleteSeat())

		r.Get("/groups", api.fetchGroups())
		r.Post("/groups", api.createGroup())

		r.Get("/ticket", api.fetchTickets())
		r.Get("/ticket/{groupId}", api.fetchGroupTicket())

		r.Post("/trigger-seating", api.triggerSeating())

		r.Get("/healthz", func(rw http.ResponseWriter, r *http.Request) {
			rw.Write([]byte("{\"ok\": true}"))
		})
	})

	api.server.Serve()

	return nil
}
