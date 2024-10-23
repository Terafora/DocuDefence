package handlers

import (
	"DocuDefense/src/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var users = make(map[string]models.User)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var userList []models.User
	for _, user := range users {
		userList = append(userList, user)
	}
	json.NewEncoder(w).Encode(userList)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	users[user.ID] = user
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var updatedUser models.User
	_ = json.NewDecoder(r.Body).Decode(&updatedUser)

	if _, ok := users[params["id"]]; ok {
		users[params["id"]] = updatedUser
		json.NewEncoder(w).Encode(updatedUser)
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if _, ok := users[params["id"]]; ok {
		delete(users, params["id"])
		json.NewEncoder(w).Encode(users)
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
	}
}
