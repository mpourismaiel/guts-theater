package api

import (
	"net/http"

	"github.com/mpourismaiel/guts-theater/seating"
)

func (a *ApiServer) triggerSeating() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		go seating.Process(*a.store.Models, a.logger)
		a.renderJSON(rw, 200, map[string]bool{"ok": true})
	}
}
