package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"mpourismaiel.dev/guts/store/models"
)

type createSectionRequest struct {
	SectionName string `json:"name"`
	Curved      bool   `json:"curved"`
	Elevation   int    `json:"elevation"`
}

type updateSectionRequest struct {
	SectionName string `json:"name"`
	Curved      bool   `json:"curved"`
	Elevation   int    `json:"elevation"`
}

func (a *ApiServer) fetchSections() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		sections, err := a.store.Models.SectionGetAll()
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		res, err := json.Marshal(sections)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(res)
	}
}

func (a *ApiServer) createSection() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var s createSectionRequest
		err := json.NewDecoder(r.Body).Decode(&s)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		section := models.Section{
			Name:      s.SectionName,
			Elevation: s.Elevation,
		}
		err = a.store.Models.SectionSave(&section)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		res, err := json.Marshal(section)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(res)
	}
}

func (a *ApiServer) updateSection() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var s updateSectionRequest
		err := json.NewDecoder(r.Body).Decode(&s)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		section, err := a.store.Models.SectionGetByName(chi.URLParam(r, "section"))
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		section.Name = s.SectionName
		section.Elevation = s.Elevation
		err = a.store.Models.SectionUpdate(section)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		res, err := json.Marshal(section)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(res)
	}
}

func (a *ApiServer) deleteSection() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		section, err := a.store.Models.SectionGetByName(chi.URLParam(r, "section"))
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		err = a.store.Models.SectionDelete(section)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		res, err := json.Marshal(section)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(res)
	}
}