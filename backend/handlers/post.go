package handlers

import (
	"encoding/json"
	"net/http"

	"real-time-forum/backend/database"
	"real-time-forum/backend/database/structs"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Check for method
	switch r.Method {
	case "GET":
		var posts []structs.Post
		var err error

		// Check for a passed search parameter
		param := r.URL.Query().Get("param")
		if param == "" {
			// If not passed, get all posts
			posts, err = database.FindAllPosts("database.db")
			if err != nil {
				http.Error(w, "500 internal server error.", http.StatusInternalServerError)
				return
			}
		} else {
			// If passed, check for the search data
			data := r.URL.Query().Get("data")
			// If not passed, return a 400 bad request
			if data == "" {
				http.Error(w, "400 bad request.", http.StatusBadRequest)
				return
			}

			// Get all posts that match the search parameter and data
			posts, err = database.FindPostByParam("database.db", param, data)
			if err != nil {
				http.Error(w, "500 internal server error.", http.StatusInternalServerError)
				return
			}
		}

		// Return the posts as JSON
		resp, err := json.Marshal(posts)
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	case "POST":
		var newPost structs.Post

		// Decode the request body into the post struct
		err := json.NewDecoder(r.Body).Decode(&newPost)
		if err != nil {
			http.Error(w, "400 bad request.", http.StatusBadRequest)
			return
		}

		// Check for a session cookie
		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "401 unauthorized.", http.StatusUnauthorized)
			return
		}

		// Get the current user from the session cookie
		user, err := database.CurrentUser("database.db", cookie.Value)
		if err != nil {
			return
		}

		// Add the new post to the database
		err = database.NewPost("database.db", newPost, user)
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Successfully added post."))
	default:
		http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}
}
