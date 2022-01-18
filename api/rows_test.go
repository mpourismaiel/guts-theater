package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mpourismaiel/guts-theater/store/models"
	"github.com/stretchr/testify/require"
)

func TestFetchRows(t *testing.T) {
	api := createApi()
	createRequestJson(api, "POST", "/section/hall/row", models.Row{
		Name: "1",
	})
	response := createRequest(api, "GET", "/section/hall/rows", nil)
	checkResponseCode(t, http.StatusOK, response.Code)

	v, _ := api.store.Models.RowGetBySection("hall")
	if len(v) == 0 {
		t.Error("Value not found in database")
	} else {
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(v)
		require.Equal(t, b.String(), response.Body.String())
	}
}

func TestCreateRow(t *testing.T) {
	api := createApi()

	response := createRequestJson(api, "POST", "/section/hall/row", models.Row{
		Name: "1",
	})

	checkResponseCode(t, http.StatusOK, response.Code)

	v, _ := api.store.Models.RowGetBySection("hall")
	if len(v) == 0 {
		t.Error("Value not found in database")
	} else {
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(v[0])
		require.Equal(t, b.String(), response.Body.String())
	}
}

func TestDeleteRow(t *testing.T) {
	api := createApi()

	createRequestJson(api, "POST", "/section/hall/row", models.Row{
		Name: "1",
	})
	response := createRequest(api, "DELETE", "/section/hall/row/1", nil)

	checkResponseCode(t, http.StatusOK, response.Code)

	v, _ := api.store.Models.RowGetBySection("hall")
	if len(v) != 0 {
		t.Error("Value not found in database")
	}
}
