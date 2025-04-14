package database

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"real-time-forum/backend/database/structs"
)

func NewPost(path string, post structs.Post, user structs.User) error {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	defer db.Close()

	date := time.Now().Format("01-02-2006 15:04:05")

	// Insert the new post into the database
	_, err = db.Exec(
		AddPost, user.Id,
		post.Category,
		post.Title,
		post.Content,
		date)
	if err != nil {
		return err
	}

	return nil
}

// Convert database rows to a slice of post structs
func ConvertRowToPost(rows *sql.Rows) ([]structs.Post, error) {
	var posts []structs.Post

	// Iterate through the query results
	for rows.Next() {
		var post structs.Post
		var date string

		// Scan the query results into the post struct
		err := rows.Scan(
			&post.Id,
			&post.User_id,
			&post.Category,
			&post.Title,
			&post.Content,
			&date)
		if err != nil {
			break
		}

		post.Date, err = time.Parse("01-02-2006 15:04:05", date)
		if err != nil {
			break
		}

		posts = append(posts, post)

	}

	// Return an empty slice with nil error when no posts are found instead of an error
	if len(posts) == 0 {
		return []structs.Post{}, nil
	}

	return posts, nil
}

func FindAllPosts(path string) ([]structs.Post, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return []structs.Post{}, errors.New("failed to open database")
	}

	defer db.Close()

	// Query the database for all posts (const GetAllPost)
	rows, err := db.Query(GetAllPost)
	if err != nil {
		return []structs.Post{}, errors.New("failed to find posts")
	}

	// Convert the rows to a slice of post structs
	posts, err := ConvertRowToPost(rows)
	if err != nil {
		return []structs.Post{}, errors.New("failed to convert rows to posts")
	}

	return posts, nil
}

// Get all posts by specific parameter
func FindPostByParam(path, param, data string) ([]structs.Post, error) {
	var query *sql.Rows

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return []structs.Post{}, errors.New("failed to open database")
	}

	defer db.Close()

	// Find post based on:
	switch param {
	case "id":
		// Convert the passed data to an integer
		id, err := strconv.Atoi(data)
		if err != nil {
			return []structs.Post{}, errors.New("id must be an integer")
		}

		query, err = db.Query(GetPostById, id)
		if err != nil {
			return []structs.Post{}, errors.New("could not find id")
		}
	case "user_id":
		query, err = db.Query(GetAllPostByUser, data)
		if err != nil {
			return []structs.Post{}, errors.New("could not find any posts by that user")
		}
	case "category":
		query, err = db.Query(GetAllPostByCategory, data)
		if err != nil {
			return []structs.Post{}, errors.New("could not find any posts with that category")
		}
	default:
		// Return an error if the parameter is invalid
		return []structs.Post{}, errors.New("invalid parameter")
	}

	// Convert the query results to a slice of post structs
	posts, err := ConvertRowToPost(query)
	if err != nil {
		return []structs.Post{}, errors.New("failed to convert")
	}

	return posts, nil
}
