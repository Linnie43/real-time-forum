package main

import (
	"log"
	"net/http"
	"real-time-forum/backend/chat"
	"real-time-forum/backend/database"
	"real-time-forum/backend/handlers"
)

func main() {
	// Initialize the database
	db := database.InitDB()
	defer db.Close()

	hub := chat.NewHub()
	go hub.Run()

	// Setup routes
	http.Handle("/styles.css", http.FileServer(http.Dir("./frontend")))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./frontend/js"))))
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/favicon.ico")
	})

	// Register routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./frontend/index.html")
	})
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/session", handlers.SessionHandler)
	http.HandleFunc("/user", handlers.UserHandler)
	http.HandleFunc("/post", handlers.PostHandler)
	http.HandleFunc("/comment", handlers.CommentHandler)
	http.HandleFunc("/message", handlers.MessageHandler)
	http.HandleFunc("/chat", handlers.ChatHandler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chat.ServeWs(hub, w, r)
	})

	// Start the server
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
