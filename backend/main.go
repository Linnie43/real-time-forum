package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	initDB()

	// Serve static files (HTML, CSS, JS)
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	// API routes
	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/logout", LogoutHandler)
	http.HandleFunc("/protected", ProtectedExampleHandler) // Example of a protected route
	http.HandleFunc("/create-post", CreatePostHandler)
	http.HandleFunc("/posts", GetPostsHandler)

	// Start the server
	port := ":8080"
	fmt.Println("Server running on http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
