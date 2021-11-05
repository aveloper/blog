package users

import (
	"net/http"

	"github.com/gorilla/mux"
)

//UserRoutes adds users routes to router
func UserRoutes(r *mux.Router, uh *Handler) {
	r = r.PathPrefix("/users").Subrouter()

	r.HandleFunc("/edit/password", uh.updateUserPassword()).Methods(http.MethodPost)
	r.HandleFunc("/edit", uh.updateUser()).Methods(http.MethodPost)
	//FIXME: it can also be changed to /delete/{id} for better use
	r.HandleFunc("/delete", uh.DeleteUser()).Methods(http.MethodPost)
	r.HandleFunc("/all", uh.getAllUser()).Methods(http.MethodGet)
	r.HandleFunc("/{id}",  uh.getUser()).Methods(http.MethodGet)
	r.HandleFunc("/", uh.addUser()).Methods(http.MethodPost)
}
