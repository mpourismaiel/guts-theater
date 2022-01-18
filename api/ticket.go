package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *ApiServer) fetchTickets() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		tickets, err := a.store.Models.TicketGetAll()
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		groupsBySeats := make(map[string][]string)
		for i, t := range tickets {
			for _, s := range t.Seats {
				groupsBySeats[s] = []string{t.GroupId, fmt.Sprint(i)}
			}
		}

		a.renderJSON(rw, 200, groupsBySeats)
	}
}

func (a *ApiServer) fetchGroupTicket() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		group, err := a.store.Models.TicketGetByGroupId(chi.URLParam(r, "groupId"))
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		a.renderJSON(rw, 200, group)
	}
}
