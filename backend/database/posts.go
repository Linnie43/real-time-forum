package database

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"real-time-forum/backend/database/structs"
)

// Attempt to insert a new post into the database
func NewPost(path string, post structs.Post, user structs.User) error {
	// Open the database
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	defer db.Close()

	date := time.Now().Format("01-02-2006 15:04:05")

	// Execute the insert statement
	_, err = db.Exec(AddPost, user.Id, post.Category, post.Title, post.Content, date)
	if err != nil {
		return err
	}

	return nil
}

// Convert post table query results to a slice of post structs
func ConvertRowToPost(rows *sql.Rows) ([]structs.Post, error) {
	var posts []structs.Post

	// Iterate through the query results
	for rows.Next() {
		var post structs.Post
		var date string

		// Scan the query results into the post struct
		err := rows.Scan(&post.Id, &post.User_id, &post.Category, &post.Title, &post.Content, &date)
		if err != nil {
			break
		}

		// Parse the date string into a time.Time struct
		post.Date, err = time.Parse("01-02-2006 15:04:05", date)
		if err != nil {
			break
		}

		// Append the post to the posts slice
		posts = append(posts, post)

	}

	// Return an error if no posts were found
	if len(posts) == 0 {
		return []structs.Post{}, sql.ErrNoRows
	}

	return posts, nil
}

// Get all posts from the database
func FindAllPosts(path string) ([]structs.Post, error) {
	// Open the database
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return []structs.Post{}, errors.New("failed to open database")
	}

	defer db.Close()

	// Find all posts
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

// Get all posts that match the passed parameter (id, user_id, cateogry)
func FindPostByParam(path, param, data string) ([]structs.Post, error) {
	var query *sql.Rows

	// Open the database
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return []structs.Post{}, errors.New("failed to open database")
	}

	defer db.Close()

	// Find posts depending on the passed parameter
	switch param {
	case "id":
		// Convert the passed data to an integer
		id, err := strconv.Atoi(data)
		if err != nil {
			return []structs.Post{}, errors.New("id must be an integer")
		}

		// Find posts by id
		query, err = db.Query(GetPostById, id)
		if err != nil {
			return []structs.Post{}, errors.New("could not find id")
		}
	case "user_id":
		// Search the database by user_id
		query, err = db.Query(GetAllPostByUser, data)
		if err != nil {
			return []structs.Post{}, errors.New("could not find any posts by that user")
		}
	case "category":
		// Search the database by category
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
