package main

import (
	"log"
	"net/http"
	"real-time-forum/backend/database"
	"real-time-forum/backend/handlers"
)

func main() {
	// Initialize the database
	db := database.InitDB()
	defer db.Close()

	// Setup the routes
	http.Handle("/styles.css", http.FileServer(http.Dir("./frontend")))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./frontend/js"))))

	// Register routes
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/session", handlers.SessionHandler)
	http.HandleFunc("/user", handlers.UserHandler)
	http.HandleFunc("/post", handlers.PostHandler)
	http.HandleFunc("/comment", handlers.CommentHandler)
	// http.HandleFunc("/message", handlers.MessageHandler)
	// mux.HandleFunc("/chat", handlers.ChatHandler)

	// Start the server
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
