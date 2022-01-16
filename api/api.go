package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"mpourismaiel.dev/guts/server"
	"mpourismaiel.dev/guts/store"
)

type ApiServer struct {
	server *server.Server
	store  *store.Orm
}

func New(port string) error {
	s, err := server.New(port)
	if err != nil {
		return fmt.Errorf("failed to create server for API: %v", err)
	}

	o := store.New("guts")

	api := ApiServer{
		server: s,
		store:  o,
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
		r.Put("/section/{section}/row/{row}", api.updateRow())
		r.Delete("/section/{section}/row/{row}", api.deleteRow())

		r.Get("/seats", api.fetchSeats())

		r.Get("/section/{section}/seats", api.fetchSeatsBySection())
		r.Post("/section/{section}/row/{row}/seat", api.createSeats())
		r.Put("/section/{section}/row/{row}/seat/{seat}", api.updateSeat())
		r.Delete("/section/{section}/row/{row}/seat/{seat}", api.deleteSeat())

		r.Post("/seating-trigger", api.triggerSeating())

		r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("{\"ok\": true}"))
		})
	})

	api.server.Serve()

	return nil
}
