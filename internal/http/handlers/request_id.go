package handlers

import (
	"github.com/aveloper/blog/internal/blogcontext"
	"github.com/google/uuid"
	"net/http"
)

const (
	// httpHeaderNameRequestID has the name of the header for request ID
	httpHeaderNameRequestID = "X-Request-ID"
)

// AssignRequestIDHandler is handler to assign request ID to each incoming request
func AssignRequestIDHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		ctx := blogcontext.AddRequestID(r.Context(), id)

		w.Header().Set(httpHeaderNameRequestID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
