package backend

import (
	"encoding/json"
	"log"
	"net/http"
)

// VerifySession checks if the session ID exists in the database
func VerifySession(r *http.Request) (bool, int, string) {
	var userID int
	var username string
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false, 0, ""
	}

	err = db.QueryRow("SELECT user_id FROM Session WHERE id = ? AND status = 'active'", cookie.Value).Scan(&userID)
	if err != nil {
		log.Println("No userID found for the cookie")
		return false, 0, ""
	}

	err = db.QueryRow("SELECT username FROM User WHERE id = ?", userID).Scan(&username)
	if err != nil {
		log.Println("No username found")
		return false, 0, ""
	}

	return true, userID, username
}

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	loggedIn, userID, _ := VerifySession(r)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"loggedIn": loggedIn,
		"userID":   userID,
	})
}
