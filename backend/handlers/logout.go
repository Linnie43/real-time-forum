package handlers

import (
	"database/sql"
	"net/http"

	"real-time-forum/backend/database"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Prevent the endpoint from being accessed by other URL paths
	if r.URL.Path != "/logout" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Only allow POST requests
	if r.Method != "POST" {
		http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Open the database connection
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		return
	}

	defer db.Close()

	// Check if user has a session cookie
	cookie, err := r.Cookie("session")
	if err != nil {
		return
	}

	// Delete the session cookie from the database
	_, err = db.Exec(database.RemoveCookie, cookie.Value)
	if err != nil {
		http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		return
	}

	// Delete the session cookie from the client
	cookie.MaxAge = -1
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully logged out."))
}
