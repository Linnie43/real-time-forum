package handlers

import (
	"encoding/json"
	"net/http"

	"real-time-forum/backend/database"
	"real-time-forum/backend/database/structs"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	// Check whether the request is a GET or POST
	switch r.Method {
	case "GET":
		// Check for a passed search parameter and data
		param := r.URL.Query().Get("param")
		data := r.URL.Query().Get("data")
		if param == "" || data == "" {
			http.Error(w, "400 bad request.", http.StatusBadRequest)
			return
		}

		// Find the comment(s) based on the passed search parameter and data
		comments, err := database.FindCommentByParam("database.db", param, data)
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}

		// Return the comment(s) as JSON
		resp, err := json.Marshal(comments)
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	case "POST":
		var newComment structs.Comment

		// Decode the request body into a new comment
		err := json.NewDecoder(r.Body).Decode(&newComment)
		if err != nil {
			http.Error(w, "400 bad request.", http.StatusBadRequest)
			return
		}

		// Insert the new comment into the database
		err = database.NewComment("database.db", newComment)
		if err != nil {
			http.Error(w, "500 internal server error.", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Comment successfully created."))
	default:
		http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}
}
