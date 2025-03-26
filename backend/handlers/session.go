package handlers

import (
	"encoding/json"
	"net/http"

	"real-time-forum/backend/database"
	"real-time-forum/backend/database/structs"
)

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	// Prevent the endpoint from being accessed by other URL paths
	if r.URL.Path != "/session" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Only allow GET requests
	if r.Method != "GET" {
		http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Check if user has a session cookie
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "401 unauthorized.", http.StatusUnauthorized)
		return
	}

	// Store the user data
	var user structs.User

	// Attempt getting the user from the session cookie
	user, err = database.CurrentUser("database.db", cookie.Value)
	if err != nil {
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
		http.Error(w, "400 bad request.", http.StatusBadRequest)
		return
	}

	// Send the user data as json to the client
	resp, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
