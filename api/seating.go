package api

import (
	"net/http"

	"mpourismaiel.dev/guts/seating"
)

func (a *ApiServer) triggerSeating() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		go seating.Process(*a.store.Models)
		rw.Write([]byte("{\"ok\": true}"))
	}
}
