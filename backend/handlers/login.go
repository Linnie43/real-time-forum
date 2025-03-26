package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/mail"
	"strconv"

	"real-time-forum/backend/database"
	"real-time-forum/backend/database/structs"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Prevent the endpoint from being accessed by other URL paths
	if r.URL.Path != "/login" {
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

	// Store the unmarshalled login data
	var login structs.Login

	// Decode the request body into the login struct
	err = json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, "400 bad request.", http.StatusBadRequest)
		return
	}

	// Username type is either email or username
	var param string

	// Check if the login entry is an email
	if _, err := mail.ParseAddress(login.Entry); err != nil {
		param = "username"
	} else {
		param = "email"
	}

	// Search for the user in the database
	user, err := database.GetUser("database.db", param, login.Entry)
	if err != nil {
		http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		return
	}

	// Compare the password hash with the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		http.Error(w, "401 unauthorized.", http.StatusUnauthorized)
		return
	}

	// Remove expired cookie based on valid user login
	_, err = db.Exec(database.RemoveCookie, user.Id)
	if err != nil {
		http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		return
	}

	// Check for session cookie, and generate a new one if it doesn't exist
	cookie, err := r.Cookie("session")
	if err != nil {
		// Generate the session uuid
		sessionId, err := uuid.NewV4()
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}

		// Create and set the cookie
		cookie = &http.Cookie{
			Name:     "session",
			Value:    sessionId.String(),
			HttpOnly: true,
			Path:     "/",
			MaxAge:   60 * 60 * 24, // 1 day
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		}
		http.SetCookie(w, cookie)
	}

	// Insert the cookie into the database
	_, err = db.Exec(database.AddSession, cookie.Value, user.Id)
	if err != nil {
		http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully logged in as " + user.Username + "(" + strconv.Itoa(user.Id) + ")."))
}
