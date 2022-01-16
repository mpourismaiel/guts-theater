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
		r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("{\"ok\": true}"))
		})
	})

	api.server.Serve()

	return nil
}
