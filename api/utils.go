package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mpourismaiel/guts-theater/config"
	"go.uber.org/zap"
)

func createApi() *ApiServer {
	api, _ := New(&config.Config{
		Address:    "localhost",
		Port:       "4000",
		DbHost:     "localhost",
		DbUser:     "admin",
		DbPassword: "password",
		DbName:     fmt.Sprintf("test-%s", randStringRunes(8)),
		TestMode:   true,
	}, zap.NewExample())

	return api
}

// createRequest creates a simple http request
func createRequest(api *ApiServer, method string, path string, v []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(v))
	response := executeRequest(req, api)
	return response
}

// createRequestJson create a json http request
func createRequestJson(api *ApiServer, method string, path string, v interface{}) *httptest.ResponseRecorder {
	body, _ := json.Marshal(v)
	return createRequest(api, method, path, body)
}

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

func init() {
	rand.Seed(time.Now().UnixNano())
}

// randStringRunes is a simple util function to generate a random string
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
