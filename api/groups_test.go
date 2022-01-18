package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// executeRequest, creates a new ResponseRecorder
// then executes the request by calling ServeHTTP in the router
// after which the handler writes the response to the response recorder
// which we can then inspect.
func executeRequest(req *http.Request, s *ApiServer) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.server.Router().ServeHTTP(rr, req)

	return rr
}

// checkResponseCode is a simple utility to check the response code
// of the response
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestFetchGroups(t *testing.T) {
	api, _ := New("localhost", "4000", "localhost", "admin", "password", zap.NewExample())

	req, _ := http.NewRequest("GET", "/groups", nil)
	response := executeRequest(req, api)

	checkResponseCode(t, http.StatusOK, response.Code)

	groups, _ := api.store.Models.GroupGetAll()
	if len(groups) == 0 {
		require.Equal(t, "[]", response.Body.String())
	} else {
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(groups)
		require.Equal(t, b.String(), response.Body.String())
	}
}

func TestCreateGroup(t *testing.T) {
	api, _ := New("localhost", "4000", "localhost", "admin", "password", zap.NewExample())

	req, _ := http.NewRequest("POST", "/groups", bytes.NewBuffer([]byte("{\"aisle\": false,\"rank\": \"silver\",\"count\": 2,\"section\": \"hall\"}")))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	checkResponseCode(t, http.StatusOK, response.StatusCode)
	body, _ := ioutil.ReadAll(response.Body)

	groups, _ := api.store.Models.GroupGetAll()
	if len(groups) == 0 {
		require.Equal(t, "[]", string(body))
	} else {
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(groups)
		require.Equal(t, b.String(), string(body))
	}
}
