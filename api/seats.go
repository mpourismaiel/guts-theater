package api

import (
	"encoding/json"
	"net/http"
)

func (a *ApiServer) fetchSeats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seats, err := a.store.Models.SeatGetAll()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		res, _ := json.MarshalIndent(seats, "", "  ")
		w.Write(res)
	}
}
