package main

import (
	"log"
	"net/http"
	"real-time-forum/backend"
	"real-time-forum/database"
	"text/template"
)

func main() {

	// Parse and serve the template
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		log.Fatal("Error parsing template")
		return
	}

	http.Handle("/assets/", http.FileServer(http.Dir(".")))
	http.Handle("/js/", http.FileServer(http.Dir(".")))

	// Initialize database
	db := database.InitDB()
	defer db.Close()

	database.MakeTables(db)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Error executing template", http.StatusInternalServerError)
		}
	})

	// One API Handler for api calls
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		backend.APIHandler(w, r, db)
	})

	log.Println("Server is running on http://localhost:8080")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
