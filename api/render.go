package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ErrResponse struct {
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

func (s *ApiServer) renderErrInternal(rw http.ResponseWriter, err error) {
	s.renderJSON(rw, http.StatusInternalServerError, ErrResponse{Status: "internal error", Error: errString(err)})
}

func (s *ApiServer) renderStringAsJSON(rw http.ResponseWriter, code int, v string) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(code)
	_, _ = rw.Write([]byte(v))
}

func (s *ApiServer) renderString(rw http.ResponseWriter, code int, v string) {
	rw.Header().Set("Content-Type", "plain/text")
	rw.WriteHeader(code)
	_, _ = rw.Write([]byte(v))
}

func (s *ApiServer) renderJSON(rw http.ResponseWriter, code int, v interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(v); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		fields := []zapcore.Field{
			zap.String("error", errString(err)),
		}
		s.logger.Error("failed to encode json", fields...)
	} else {
		rw.WriteHeader(code)
	}

	_, _ = rw.Write(b.Bytes())
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
