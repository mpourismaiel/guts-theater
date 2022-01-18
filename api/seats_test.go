package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mpourismaiel/guts-theater/store/models"
	"github.com/stretchr/testify/require"
)

func TestFetchSeats(t *testing.T) {
	api := createApi()
	createRequestJson(api, "POST", "/section/hall/row/0/seat", []models.Seat{{
		Name:   "0",
		Rank:   "silver",
		Broken: false,
		Aisle:  false,
	}})
	response := createRequest(api, "GET", "/section/hall/seats", nil)
	checkResponseCode(t, http.StatusOK, response.Code)

	v, _ := api.store.Models.SeatGetBySection("hall")
	if len(v) == 0 {
		t.Error("Value not found in database")
	} else {
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(v)
		require.Equal(t, b.String(), response.Body.String())
	}
}

func TestCreateSeat(t *testing.T) {
	api := createApi()

	response := createRequestJson(api, "POST", "/section/hall/row/0/seat", []models.Seat{{
		Name:   "0",
		Rank:   "silver",
		Broken: false,
		Aisle:  false,
	}})

	checkResponseCode(t, http.StatusOK, response.Code)

	v, _ := api.store.Models.SeatGetBySection("hall")
	if len(v) == 0 {
		t.Error("Value not found in database")
	} else {
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(v)
		require.Equal(t, b.String(), response.Body.String())
	}
}

func TestUpdateSeat(t *testing.T) {
	api := createApi()

	createRequestJson(api, "POST", "/section/hall/row/0/seat", []models.Seat{{
		Name:   "0",
		Rank:   "silver",
		Broken: false,
		Aisle:  false,
	}})
	response := createRequestJson(api, "PUT", "/section/hall/row/0/seat/0", models.Seat{
		Name:   "0",
		Rank:   "gold",
		Broken: false,
		Aisle:  false,
	})

	checkResponseCode(t, http.StatusOK, response.Code)

	v, _ := api.store.Models.SeatGetBySection("hall")
	if len(v) == 0 {
		t.Error("Value not found in database")
	} else {
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(v[0])
		require.Equal(t, b.String(), response.Body.String())
	}
}

func TestDeleteSeat(t *testing.T) {
	api := createApi()

	createRequestJson(api, "POST", "/section/hall/row/0/seat", []models.Seat{{
		Name:   "0",
		Rank:   "silver",
		Broken: false,
		Aisle:  false,
	}})
	response := createRequest(api, "DELETE", "/section/hall/row/0/seat/0", nil)

	checkResponseCode(t, http.StatusOK, response.Code)

	v, _ := api.store.Models.SeatGetBySection("hall")
	if len(v) != 0 {
		t.Error("Value not found in database")
	}
}
