package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (a *ApiServer) fetchTickets() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		tickets, err := a.store.Models.TicketGetAll()
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		groupsBySeats := make(map[string][]string)
		for i, t := range tickets {
			for _, s := range t.Seats {
				groupsBySeats[s] = []string{t.GroupId, fmt.Sprint(i)}
			}
		}

		res, err := json.Marshal(groupsBySeats)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(res)
	}
}

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
