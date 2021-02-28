package service

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/login", HandleLogin).Methods("POST")
	r.HandleFunc("/register", HandleRegistration).Methods("POST")
	r.HandleFunc("/validateToken", HandleValidateToken).Methods("POST")
	return r
}
