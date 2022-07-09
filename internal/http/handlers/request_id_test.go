package handlers

import (
	"github.com/aveloper/blog/internal/blogcontext"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAssignRequestIDHandler(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, err := blogcontext.GetRequestID(r.Context())
		assert.Nil(t, err)

		assert.NotEmpty(t, requestID)
	})

	handler := AssignRequestIDHandler(nextHandler)

	r := httptest.NewRequest(http.MethodGet, "/some/endpoint", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, r)

	id := rr.Header().Get(httpHeaderNameRequestID)

	assert.NotEmpty(t, id)
}
