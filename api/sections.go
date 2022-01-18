package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mpourismaiel/guts-theater/store/models"
)

type createSectionRequest struct {
	SectionName string `json:"name"`
	Curved      bool   `json:"curved"`
	Elevation   int    `json:"elevation"`
}

type updateSectionRequest struct {
	Curved    bool `json:"curved"`
	Elevation int  `json:"elevation"`
}

func (a *ApiServer) fetchSections() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		sections, err := a.store.Models.SectionGetAll()
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		a.renderJSON(rw, 200, sections)
	}
}

func (a *ApiServer) createSection() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var s createSectionRequest
		err := json.NewDecoder(r.Body).Decode(&s)
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		section := models.Section{
			Name:      s.SectionName,
			Elevation: s.Elevation,
		}
		err = a.store.Models.SectionSave(&section)
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		a.renderJSON(rw, 200, section)
	}
}

func (a *ApiServer) updateSection() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var s updateSectionRequest
		err := json.NewDecoder(r.Body).Decode(&s)
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		section, err := a.store.Models.SectionGetByName(chi.URLParam(r, "section"))
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		section.Curved = s.Curved
		section.Elevation = s.Elevation
		err = a.store.Models.SectionUpdate(section)
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		a.renderJSON(rw, 200, section)
	}
}

func (a *ApiServer) deleteSection() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		section, err := a.store.Models.SectionGetByName(chi.URLParam(r, "section"))
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		err = a.store.Models.SectionDelete(section)
		if err != nil {
			a.renderErrInternal(rw, err)
			return
		}

		a.renderJSON(rw, 200, section)
	}
}
