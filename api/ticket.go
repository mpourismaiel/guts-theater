package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *ApiServer) fetchGroupTicket() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		group, err := a.store.Models.TicketGetByGroupId(chi.URLParam(r, "groupId"))
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		res, err := json.Marshal(group)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(res)
	}
}
