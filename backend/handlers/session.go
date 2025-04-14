package handlers

import (
	"encoding/json"
	"net/http"

	"real-time-forum/backend/database"
	"real-time-forum/backend/database/structs"
)

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/session" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Only allow GET requests
	if r.Method != "GET" {
		http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Check for the session cookie
	cookie, err := r.Cookie("session")
	if err != nil {
		http.Error(w, "401 unauthorized.", http.StatusUnauthorized)
		return
	}

	// Store user data
	var user structs.User

	// Get user from the database by session cookie
	user, err = database.CurrentUser("database.db", cookie.Value)
	if err != nil {
		cookie.MaxAge = -1 // time out the cookie
		http.SetCookie(w, cookie) // return the cookie to the client
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
