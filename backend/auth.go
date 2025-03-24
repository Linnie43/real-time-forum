package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/mattn/go-sqlite3"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("super-secret-key")) // Change this key!

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Nickname  string `json:"nickname"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Age       int    `json:"age"`
		Gender    string `json:"gender"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create a new user
	user, err := NewUser(input.Nickname, input.Email, input.Password, input.Age, input.Gender, input.FirstName, input.LastName)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Insert into database
	_, err = db.Exec(`INSERT INTO users (id, nickname, email, password, age, gender, first_name, last_name)
					  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Nickname, user.Email, user.Password, user.Age, user.Gender, user.FirstName, user.LastName)
	if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique { // Unique constraint failed
		if sqlErr, ok := err.(sqlite3.Error); ok && sqlErr.Code == 2067 { // Unique constraint failed
			http.Error(w, "Nickname or email already in use", http.StatusConflict)
			return
		}
		log.Println("DB error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

// LoginHandler handles user authentication
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		Identifier string `json:"identifier"` // Can be email or nickname
		Password   string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find user by nickname or email
	var user User
	err = db.QueryRow(`SELECT id, nickname, email, password FROM users WHERE email = ? OR nickname = ?`,
		input.Identifier, input.Identifier).
		Scan(&user.ID, &user.Nickname, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Println("DB error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Check password
	if !CheckPasswordHash(input.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create session
	session, _ := store.Get(r, "session")
	session.Values["user_id"] = user.ID
	session.Values["nickname"] = user.Nickname
	session.Save(r, w)

	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

// LogoutHandler clears the session
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Options.MaxAge = -1 // Deletes the session
	session.Save(r, w)

	json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful"})
}

// GetUserFromSession retrieves the logged-in user ID
func GetUserFromSession(r *http.Request) (string, bool) {
	session, _ := store.Get(r, "session")
	userID, ok := session.Values["user_id"].(string)
	return userID, ok
}

// ProtectedExampleHandler is an example of a protected route
func ProtectedExampleHandler(w http.ResponseWriter, r *http.Request) {
	userID, loggedIn := GetUserFromSession(r)
	if !loggedIn {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Welcome, user " + userID})
}
