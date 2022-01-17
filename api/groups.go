package api

import (
	"encoding/json"
	"net/http"

	"mpourismaiel.dev/guts/store/models"
)

type createGroupRequest struct {
	Aisle   bool   `json:"aisle"`
	Rank    string `json:"rank"`
	Count   int    `json:"count"`
	Section string `json:"section"`
}

func (a *ApiServer) fetchGroups() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		groups, err := a.store.Models.GroupGetAll()
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		if len(groups) == 0 {
			rw.Write([]byte("[]"))
			return
		}

		res, err := json.Marshal(groups)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(res)
	}
}

func (a *ApiServer) createGroup() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var group createGroupRequest
		err := json.NewDecoder(r.Body).Decode(&group)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		newGroup := models.Group{
			Aisle:   group.Aisle,
			Rank:    group.Rank,
			Count:   group.Count,
			Section: group.Section,
		}
		err = a.store.Models.GroupSave(&newGroup)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}

		res, err := json.Marshal(newGroup)
		if err != nil {
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Write(res)
	}
}
