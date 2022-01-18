package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mpourismaiel/guts-theater/store/models"
	"github.com/stretchr/testify/require"
)

func TestFetchGroups(t *testing.T) {
	api := createApi()
	createRequestJson(api, "POST", "/groups", models.Group{
		Aisle:   false,
		Rank:    "silver",
		Count:   2,
		Section: "hall",
	})
	response := createRequest(api, "GET", "/groups", nil)

	checkResponseCode(t, http.StatusOK, response.Code)

	v, _ := api.store.Models.GroupGetAll()
	if len(v) == 0 {
		t.Error("Value not found in database")
	} else {
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(v)
		require.Equal(t, b.String(), response.Body.String())
	}
}

func TestCreateGroup(t *testing.T) {
	api := createApi()

	response := createRequestJson(api, "POST", "/groups", models.Group{
		Aisle:   false,
		Rank:    "silver",
		Count:   2,
		Section: "hall",
	})

	checkResponseCode(t, http.StatusOK, response.Code)

	v, _ := api.store.Models.GroupGetAll()
	if len(v) == 0 {
		t.Error("Value not found in database")
	} else {
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(v[0])
		require.Equal(t, b.String(), response.Body.String())
	}
}
