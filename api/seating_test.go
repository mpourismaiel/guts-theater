package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTriggerSeating(t *testing.T) {
	api := createApi()
	response := createRequest(api, "POST", "/trigger-seating", nil)

	checkResponseCode(t, http.StatusOK, response.Code)
	require.Equal(t, "{\"ok\":true}\n", response.Body.String())
}
