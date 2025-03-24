package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

// CreatePostHandler handles new post creation
func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Check if the user is logged in
	userID, loggedIn := GetUserFromSession(r)
	if !loggedIn {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate post fields
	if post.Title == "" || post.Content == "" || post.Category == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Generate a new post ID
	newID, _ := uuid.NewV4()

	// Insert into database
	_, err = db.Exec(`INSERT INTO posts (id, user_id, title, content, category, created_at)
					  VALUES (?, ?, ?, ?, ?, ?)`,
		newID.String(), userID, post.Title, post.Content, post.Category, time.Now())

	if err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Post created successfully"})
}

// GetPostsHandler retrieves all posts
func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, user_id, title, content, category, created_at FROM posts ORDER BY created_at DESC`)
	if err != nil {
		log.Println("Database error:", err)
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content, &post.Category, &post.CreatedAt)
		if err != nil {
			http.Error(w, "Error scanning posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	json.NewEncoder(w).Encode(posts)
}
