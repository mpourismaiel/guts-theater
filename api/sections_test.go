package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/mpourismaiel/guts-theater/store/models"
	"github.com/stretchr/testify/require"
)

func TestFetchSections(t *testing.T) {
	api := createApi()
	createRequestJson(api, "POST", "/section", models.Section{
		Name:      "hall",
		Elevation: 0,
		Curved:    false,
	})
	response := createRequest(api, "GET", "/section", nil)

	checkResponseCode(t, http.StatusOK, response.Code)

	v, _ := api.store.Models.SectionGetAll()
	if len(v) == 0 {
		t.Error("Value not found in database")
	} else {
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(v)
		require.Equal(t, b.String(), response.Body.String())
	}
}

func TestCreateSections(t *testing.T) {
	api := createApi()
	response := createRequestJson(api, "POST", "/section", models.Section{
		Name:      "hall",
		Elevation: 0,
		Curved:    false,
	})

	checkResponseCode(t, http.StatusOK, response.Code)

	v, _ := api.store.Models.SectionGetAll()
	if len(v) == 0 {
		t.Error("Value not found in database")
	} else {
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(v[0])
		require.Equal(t, b.String(), response.Body.String())
	}
}

func TestUpdateSection(t *testing.T) {
	api := createApi()
	createRequestJson(api, "POST", "/section", models.Section{
		Name:      "hall",
		Elevation: 0,
		Curved:    false,
	})
	response := createRequestJson(api, "PUT", "/section/hall", models.Section{
		Elevation: 1,
		Curved:    false,
	})

	checkResponseCode(t, http.StatusOK, response.Code)

	v, _ := api.store.Models.SectionGetAll()
	if len(v) == 0 {
		t.Error("Value not found in database")
	} else {
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(v[0])
		require.Equal(t, b.String(), response.Body.String())
	}
}

func TestDeleteSection(t *testing.T) {
	api := createApi()
	createRequestJson(api, "POST", "/section", models.Section{
		Name:      "hall",
		Elevation: 0,
		Curved:    false,
	})
	response := createRequest(api, "DELETE", "/section/hall", nil)

	checkResponseCode(t, http.StatusOK, response.Code)

	v, _ := api.store.Models.SectionGetAll()
	if len(v) != 0 {
		t.Error("Value was not deleted from database")
	}
}
