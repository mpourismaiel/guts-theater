package api

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealthz(t *testing.T) {
	api := createApi()
	response := createRequest(api, "GET", "/healthz", nil)

	checkResponseCode(t, http.StatusOK, response.Code)
	require.Equal(t, "{\"ok\":true}\n", response.Body.String())
}
