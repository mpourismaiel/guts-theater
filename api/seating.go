package api

import (
	"net/http"

	"github.com/mpourismaiel/guts-theater/seating"
)

func (a *ApiServer) triggerSeating(syncMode bool) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if syncMode {
			seating.Process(*a.store.Models, a.logger)
		} else {
			go seating.Process(*a.store.Models, a.logger)
		}
		a.renderJSON(rw, 200, map[string]bool{"ok": true})
	}
}
