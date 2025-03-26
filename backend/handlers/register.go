package handlers

import (
	"encoding/json"
	"net/http"

	"real-time-forum/backend/database"
	"real-time-forum/backend/database/structs"

	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Prevent the endpoint from being accessed by other URL paths
	if r.URL.Path != "/register" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Only allow POST requests
	if r.Method != "POST" {
		http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}

	// Store the unmarshalled register data
	var newUser structs.User

	// Decode the request body into newUser
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "400 bad request.", http.StatusBadRequest)
		return
	}

	// Generate the password hash for the user
	passwordHash, err := GenerateHash(newUser.Password)
	if err != nil {
		http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		return
	}

	newUser.Password = passwordHash

	// Attempt adding the new user to the database
	err = database.NewUser("database.db", newUser)
	if err != nil {
		http.Error(w, "500 internal server error.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successful registration."))
}

// Generates a hash from a given password
func GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)

	return string(hash), err
}
