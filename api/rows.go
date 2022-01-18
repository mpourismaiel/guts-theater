package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"mpourismaiel.dev/guts/store/models"
)

type createRowRequest struct {
	RowName string `json:"name"`
}

func (a *ApiServer) fetchRowsBySection() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rows, err := a.store.Models.RowGetBySection(chi.URLParam(r, "section"))
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		a.renderJSON(rw, 200, rows)
	}
}

func (a *ApiServer) createRow() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var row createRowRequest
		err := json.NewDecoder(r.Body).Decode(&row)
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		newRow := models.Row{
			Name:    row.RowName,
			Section: chi.URLParam(r, "section"),
		}
		err = a.store.Models.RowSave(&newRow)
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		a.renderJSON(rw, 200, newRow)
	}
}

func (a *ApiServer) deleteRow() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		row, err := a.store.Models.RowGetByName(chi.URLParam(r, "section"), chi.URLParam(r, "row"))
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		err = a.store.Models.RowDelete(row)
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		a.renderJSON(rw, 200, row)
	}
}
