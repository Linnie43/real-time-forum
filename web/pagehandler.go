package web

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"
)

var db *sql.DB

var tmpl = template.Must(template.ParseGlob("templates/*.html"))

func PageHandler(w http.ResponseWriter, r *http.Request, database *sql.DB) {

	data := PageDetails{}

	db = database

	// trim api prefix from url

	// remove data since we dont need it

	switch r.URL.Path {
	case "/":
		HomePage(w, r)
	case "/login":
		Login(w, r)
	case "/signup":
		SignUp(w, r, &data)
	case "/create-post":
		CreatePost(w, r, &data)
	default:
		if strings.HasPrefix(r.URL.Path, "/post") {
			PostHandler(w, r, &data)
		} else if strings.HasPrefix(r.URL.Path, "/logout") {
			Logout(w, r, &data)
		} else {
			ErrorHandler(w, "Page Not Found", http.StatusNotFound)
		}
	}
}

// RenderTemplate handles the rendering of HTML templates with provided data
func RenderTemplate(w http.ResponseWriter, t string, data interface{}) {

	err := tmpl.ExecuteTemplate(w, t+".html", data)
	if err != nil {
		log.Println("Error executing template:", err)
		ErrorHandler(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// ErrorHandler handles the rendering of error pages
func ErrorHandler(w http.ResponseWriter, errorMessage string, statusCode int) {
	// send back http status code instead of render error
	// w.WriteHeader(statusCode)

	// err := tmpl.ExecuteTemplate(w, "error.html", map[string]string{
	// 	"ErrorMessage": errorMessage,
	// })
	// if err != nil {
	// 	log.Println("Error executing template error.html:", err)
	// 	http.Error(w, errorMessage, statusCode)
	// 	return
	// }
}

// VerifySession checks if the session ID exists in the database
func VerifySession(r *http.Request) (bool, int, string) {
	var userID int
	var username string
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return false, 0, ""
	}

	err = db.QueryRow("SELECT user_id FROM Session WHERE id = ?", cookie.Value).Scan(&userID)
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
