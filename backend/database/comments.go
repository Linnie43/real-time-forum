package database

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"real-time-forum/backend/database/structs"
)

func NewComment(path string, comment structs.Comment) error {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	defer db.Close()

	date := time.Now().Format("01-02-2006 15:04:05")

	// Insert the new comment into the database
	_, err = db.Exec(
		AddComment, 
		comment.Post_id, 
		comment.User_id, 
		comment.Content, 
		date)
	if err != nil {
		return err
	}

	return nil
}

// Convert database rows to a slice of comment structs
func ConvertRowToComment(rows *sql.Rows) ([]structs.Comment, error) {
	var comments []structs.Comment

	// Iterate through each row
	for rows.Next() {
		var comment structs.Comment
		var date string

		// Scan the row into the comment struct
		err := rows.Scan(
			&comment.Id, 
			&comment.Post_id, 
			&comment.User_id, 
			&comment.Content, 
			&date)
		if err != nil {
			break
		}

		// Parse the date string into a time.Time struct
		comment.Date, err = time.Parse("01-02-2006 15:04:05", date)
		if err != nil {
			break
		}

		// Append the comment to the slice
		comments = append(comments, comment)
	}

	return comments, nil
}

func FindCommentByParam(path, param, data string) ([]structs.Comment, error) {
	var query *sql.Rows

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return []structs.Comment{}, errors.New("failed to open database")
	}

	defer db.Close()

	// Convert data to an integer
	id, err := strconv.Atoi(data)
	if err != nil {
		return []structs.Comment{}, errors.New("must provide an integer")
	}

	switch param {
	case "id":
		// Search for a comment by id
		query, err = db.Query(GetCommentById, id)
		if err != nil {
			return []structs.Comment{}, errors.New("could not find id")
		}
	case "post_id":
		// Search for a comment by post_id
		query, err = db.Query(GetAllPostComment, id)
		if err != nil {
			return []structs.Comment{}, errors.New("could not find post_id")
		}
	case "user_id":
		// Search for a comment by user_id
		query, err = db.Query(GetAllUserComment, id)
		if err != nil {
			return []structs.Comment{}, errors.New("could not find user_id")
		}
	default:
		return []structs.Comment{}, errors.New("invalid search parameter")
	}

	// Convert the query results to a slice of comment structs
	comments, err := ConvertRowToComment(query)
	if err != nil {
		return []structs.Comment{}, errors.New("failed to convert")
	}

	return comments, nil
}
