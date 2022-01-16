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

type updateRowRequest struct {
	RowName string `json:"name"`
}

func (a *ApiServer) fetchRowsBySection() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rows, err := a.store.Models.RowGetBySection(chi.URLParam(r, "section"))
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		res, err := json.Marshal(rows)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(res)
	}
}

func (a *ApiServer) createRow() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var row createRowRequest
		err := json.NewDecoder(r.Body).Decode(&r)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		newRow := models.Row{
			Name:    row.RowName,
			Section: chi.URLParam(r, "section"),
		}
		err = a.store.Models.RowSave(&newRow)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		res, err := json.Marshal(newRow)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(res)
	}
}

func (a *ApiServer) updateRow() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var row updateRowRequest
		err := json.NewDecoder(r.Body).Decode(&row)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		foundRow, err := a.store.Models.RowGetByName(chi.URLParam(r, "section"), chi.URLParam(r, "row"))
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		foundRow.Name = row.RowName
		err = a.store.Models.RowUpdate(foundRow)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		res, err := json.Marshal(foundRow)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(res)
	}
}

func (a *ApiServer) deleteRow() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		row, err := a.store.Models.RowGetByName(chi.URLParam(r, "section"), chi.URLParam(r, "row"))
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		err = a.store.Models.RowDelete(row)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		res, err := json.Marshal(row)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(res)
	}
}
