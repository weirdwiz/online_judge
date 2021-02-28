package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/weirdwiz/online_judge/authentication/internal/app/dbclient"
	"github.com/weirdwiz/online_judge/authentication/internal/app/model"
)

func HandleRegistration(w http.ResponseWriter, r *http.Request) {
	var user model.User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("Error Decoding User"))
	}
	success, err := DBClient.CreateUser(user)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	fmt.Fprintf(w, "Status: %t", success)
}

func WriteError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, err.Error())
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		WriteError(w, http.StatusBadRequest, fmt.Errorf("Error Decoding User"))
	}
	status, err := DBClient.Login(user.Email, user.Password)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}
	fmt.Fprintf(w, status)
}

func HandleValidateToken(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Not Yet Implemented")
}

var DBClient dbclient.IDBClient
