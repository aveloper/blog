package users

import (
	"net/http"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router, uh *Handler) {
	r = r.PathPrefix("/users").Subrouter()

	r.HandleFunc("/", uh.addUser()).Methods(http.MethodPost)
}
