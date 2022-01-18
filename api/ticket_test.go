package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/mpourismaiel/guts-theater/store/models"
	"github.com/stretchr/testify/require"
)

func TestFetchTickets(t *testing.T) {
	api := createApi()
	response := createRequest(api, "GET", "/ticket", nil)

	checkResponseCode(t, http.StatusOK, response.Code)

	v, _ := api.store.Models.TicketGetAll()
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(v)
	require.Equal(t, "{}\n", response.Body.String())
}

func TestFetchTicketsWhenHasGroup(t *testing.T) {
	api := createApi()
	createRequestJson(api, "POST", "/groups", models.Group{
		Aisle:   false,
		Rank:    "silver",
		Count:   2,
		Section: "hall",
	})
	createRequestJson(api, "POST", "/section", models.Section{
		Name:      "hall",
		Elevation: 0,
		Curved:    false,
	})
	createRequestJson(api, "POST", "/section/hall/row", models.Row{
		Name: "0",
	})
	createRequestJson(api, "POST", "/section/hall/row/0/seat", []models.Seat{{
		Name:   "0",
		Rank:   "silver",
		Broken: false,
		Aisle:  false,
	}})
	createRequestJson(api, "POST", "/section/hall/row/0/seat", []models.Seat{{
		Name:   "2",
		Rank:   "silver",
		Broken: false,
		Aisle:  false,
	}})
	createRequest(api, "POST", "/trigger-seating/sync", nil)
	response := createRequest(api, "GET", "/ticket", nil)

	checkResponseCode(t, http.StatusOK, response.Code)

	v, _ := api.store.Models.TicketGetAll()
	groupsBySeats := make(map[string][]string)
	for i, t := range v {
		for _, s := range t.Seats {
			groupsBySeats[s] = []string{t.GroupId, fmt.Sprint(i)}
		}
	}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(groupsBySeats)
	require.Equal(t, b.String(), response.Body.String())
}
