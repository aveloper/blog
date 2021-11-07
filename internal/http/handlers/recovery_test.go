package handlers

import (
	"github.com/aveloper/blog/internal/http/response"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecoveryHandler_web(t *testing.T) {
	log := zaptest.NewLogger(t)
	jw := response.NewJSONWriter(log)

	nextHandler := http.HandlerFunc(handlerThatPanics)

	handler := RecoveryHandler(log, jw)(nextHandler)

	r := httptest.NewRequest(http.MethodGet, "/some/endpoint", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.NotEqual(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestRecoveryHandler_API(t *testing.T) {
	log := zaptest.NewLogger(t)
	jw := response.NewJSONWriter(log)

	nextHandler := http.HandlerFunc(handlerThatPanics)

	handler := RecoveryHandler(log, jw)(nextHandler)

	r := httptest.NewRequest(http.MethodGet, "/api/endpoint", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
}

func TestRecoveryHandler_noPanic(t *testing.T) {
	log := zaptest.NewLogger(t)
	jw := response.NewJSONWriter(log)

	nextHandler := http.HandlerFunc(normalHandler)

	handler := RecoveryHandler(log, jw)(nextHandler)

	r := httptest.NewRequest(http.MethodGet, "/api/endpoint", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, r)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func handlerThatPanics(_ http.ResponseWriter, _ *http.Request) {
	panic("test panic")
}

func normalHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
