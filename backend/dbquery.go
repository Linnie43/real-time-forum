package backend

import (
	"log"
	"real-time-forum/backend/database"
	"real-time-forum/backend/database/structs"
	"strings"
)

// GetCategories retrieves all categories from the database
func GetCategories() ([]structs.CategoryDetails, error) {
	rows, err := db.Query(database.GetAllPostByCategory) // Use the correct query constant
	if err != nil {
		log.Println("Error retrieving categories:", err)
		return nil, err
	}
	defer rows.Close()

	var categories []structs.CategoryDetails
	for rows.Next() {
		var category structs.CategoryDetails
		if err := rows.Scan(&category.CategoryID, &category.CategoryName); err != nil {
			log.Println("Error scanning category:", err)
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

// GetPostDetails fetches the details of a specific post from the database
func GetPostDetails(postID, userID int) (*structs.PostDetails, error) {
	row := db.QueryRow(database.GetPostById, postID) // Use the correct query constant
	var err error
	post := structs.PostDetails{}
	var categories string
	err = row.Scan(
		&post.PostID,
		&post.UserID,
		&post.Username,
		&post.PostTitle,
		&post.PostContent,
		&post.CreatedAt,
		&categories,
	)

	if err != nil {
		log.Println("Error scanning rows")
		return nil, err
	}

	if categories != "" {
		post.Categories = strings.Split(categories, ",")
	}

	return &post, nil
}
