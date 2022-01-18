package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"mpourismaiel.dev/guts/seating"
	"mpourismaiel.dev/guts/store/models"
)

type createSeatRequest struct {
	Name   string `json:"name"`
	Rank   string `json:"rank"`
	Broken bool   `json:"broken"`
	Aisle  bool   `json:"aisle"`
}

type updateSeatRequest struct {
	Name   string `json:"name"`
	Rank   string `json:"rank"`
	Broken bool   `json:"broken"`
	Aisle  bool   `json:"aisle"`
}

func (a *ApiServer) fetchSeats() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		sections, err := seating.GetSections(*a.store.Models, a.logger)
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		a.renderJSON(rw, 200, sections)
	}
}

func (a *ApiServer) fetchSeatsBySection() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		seats, err := a.store.Models.SeatGetBySection(chi.URLParam(r, "section"))
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		a.renderJSON(rw, 200, seats)
	}
}

func (a *ApiServer) createSeats() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var s []createSeatRequest
		err := json.NewDecoder(r.Body).Decode(&s)
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		section := chi.URLParam(r, "section")
		row := chi.URLParam(r, "row")

		var result []*models.Seat
		for _, seat := range s {
			newSeat := models.Seat{
				Name:    seat.Name,
				Rank:    seat.Rank,
				Broken:  seat.Broken,
				Aisle:   seat.Aisle,
				Row:     row,
				Section: section,
			}
			err = a.store.Models.SeatSave(&newSeat)
			if err != nil {
				a.renderErrInternal(rw, err)
				return
			}
			result = append(result, &newSeat)
		}

		a.renderJSON(rw, 200, result)
	}
}

func (a *ApiServer) updateSeat() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var s updateSeatRequest
		err := json.NewDecoder(r.Body).Decode(&s)
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		seat, err := a.store.Models.SeatGetByName(chi.URLParam(r, "section"), chi.URLParam(r, "row"), chi.URLParam(r, "seat"))
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		seat.Name = s.Name
		seat.Rank = s.Rank
		seat.Aisle = s.Aisle
		seat.Broken = s.Broken
		err = a.store.Models.SeatUpdate(seat)
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		a.renderJSON(rw, 200, seat)
	}
}

func (a *ApiServer) deleteSeat() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		seat, err := a.store.Models.SeatGetByName(chi.URLParam(r, "section"), chi.URLParam(r, "row"), chi.URLParam(r, "seat"))
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		err = a.store.Models.SeatDelete(seat)
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		a.renderJSON(rw, 200, seat)
	}
}
