package users

import (
	"net/http"

	"github.com/gorilla/mux"
)

//UserRoutes adds users routes to router
func UserRoutes(r *mux.Router, uh *Handler) {
	r = r.PathPrefix("/users").Subrouter()

	r.HandleFunc("/", uh.addUser()).Methods(http.MethodPost)
	r.HandleFunc("/edit", uh.updateUser()).Methods(http.MethodPost)
	r.HandleFunc("/edit/name", uh.updateUserName()).Methods(http.MethodPost)
	r.HandleFunc("/edit/role", uh.updateUserRole()).Methods(http.MethodPost)
	r.HandleFunc("/edit/email", uh.updateUserEmail()).Methods(http.MethodPost)
	r.HandleFunc("/edit/password", uh.updateUserPassword()).Methods(http.MethodPost)
	//FIXME: it can also be changed to /delete/{id} for better use
 	r.HandleFunc("/delete", uh.DeleteUser()).Methods(http.MethodPost)
}
