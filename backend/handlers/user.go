package handlers

import (
	"encoding/json"
	"net/http"

	"real-time-forum/backend/database"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != "GET" {
		http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Get id from the URL
	id := r.URL.Query().Get("id")
	if id == "" { 
		// if no id is passed, return all users
		users, err := database.FindAllUsers("database.db")
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}

		// Return the users as JSON
		resp, err := json.Marshal(users)
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	} else {
		// Find the user by the passed id
		user, err := database.FindUserByParam("database.db", "id", id)
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}

		// Return the user as JSON
		resp, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}
