package handlers

import (
	"DocuDefense/src/models"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

var users = make(map[string]models.User)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	users[user.ID] = user
	json.NewEncoder(w).Encode(user)
}
