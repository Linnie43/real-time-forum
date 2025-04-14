package database

import (
	"database/sql"
	"errors"
	"strconv"

	"real-time-forum/backend/database/structs"

	_ "github.com/mattn/go-sqlite3"
)

// To add a new user to the database
func NewUser(path string, user structs.User) error {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	defer db.Close()

	// Insert new data
	_, err = db.Exec(
		AddUser, 
		user.Username, 
		user.Email, 
		user.Firstname, 
		user.Lastname, 
		user.Gender, 
		user.DOB, 
		user.Password)
	if err != nil {
		return err
	}

	return nil
}

// Get user by parameter
func GetUser(path string, parameter string, data string) (structs.User, error) {
	var query *sql.Rows

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return structs.User{}, errors.New("failed to open the database")
	}

	defer db.Close()

	// Find user based on:
	switch parameter {
	case "id":
		// Convert the data to an integer
		id, err := strconv.Atoi(data)
		if err != nil {
			return structs.User{}, errors.New("ID must be an integer")
		}

		query, err = db.Query(GetUserById, id)
		if err != nil {
			return structs.User{}, errors.New("could not find ID")
		}
	case "username":
		query, err = db.Query(GetUserByUsername, data)
		if err != nil {
			return structs.User{}, errors.New("could not find username")
		}
	case "email":
		query, err = db.Query(GetUserByEmail, data)
		if err != nil {
			return structs.User{}, errors.New("could not find email")
		}
	default:
		return structs.User{}, errors.New("invalid parameter")
	}

	// Convert the database row into a users struct
	users, err := ConvertRowToUser(query)
	if err != nil {
		return structs.User{}, errors.New("failed to convert")
	}

	if len(users) == 0 {
		return structs.User{}, errors.New("no users found")
	}

	return users[0], nil
}
// Convert the database row into a user struct
func ConvertRowToUser(rows *sql.Rows) ([]structs.User, error) {
	var users []structs.User

	for rows.Next() {
		var user structs.User

		// Store the row data in the temporary user struct
		err := rows.Scan(
			&user.Id, 
			&user.Username, 
			&user.Email, 
			&user.Firstname, 
			&user.Lastname, 
			&user.Gender, 
			&user.DOB, 
			&user.Password)
		if err != nil {
			break
		}

		// Append the temporary user struct to the users slice
		users = append(users, user)
	}

	return users, nil
}

// Finds the currently logged in user from the cookie
func CurrentUser(path, val string) (structs.User, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return structs.User{}, err
	}

	defer db.Close()
	// Query for user based on the cookie value
	query, err := db.Query(GetSessionUser, val)
	if err != nil {
		return structs.User{}, err
	}
	// Convert the database row into a user struct
	users, err := ConvertRowToUser(query)
	if err != nil {
		return structs.User{}, err
	}

	if len(users) == 0 {
		return structs.User{}, errors.New("no users found")
	}

	return users[0], nil
}

// Finds all users in the database
func FindAllUsers(path string) ([]structs.User, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return []structs.User{}, errors.New("failed to open the database")
	}

	defer db.Close()

	// Query all users from the database
	rows, err := db.Query(GetAllUser)
	if err != nil {
		return []structs.User{}, errors.New("failed to find users")
	}

	// Convert the rows into an array of users
	users, err := ConvertRowToUser(rows)
	if err != nil {
		return []structs.User{}, errors.New("failed to convert")
	}

	return users, nil
}

// Find user from the database based on the passed parameter (id, username, email)
func FindUserByParam(path string, parameter string, data string) (structs.User, error) {
	var query *sql.Rows

	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return structs.User{}, errors.New("failed to open the database")
	}

	defer db.Close()

	// Find user by:
	switch parameter {
	case "id":
		// Convert the data to an integer
		id, err := strconv.Atoi(data)
		if err != nil {
			return structs.User{}, errors.New("ID must be an integer")
		}

		query, err = db.Query(GetUserById, id)
		if err != nil {
			return structs.User{}, errors.New("could not find ID")
		}
	case "username":
		query, err = db.Query(GetUserByUsername, data)
		if err != nil {
			return structs.User{}, errors.New("could not find username")
		}
	case "email":
		query, err = db.Query(GetUserByEmail, data)
		if err != nil {
			return structs.User{}, errors.New("could not find email")
		}
	default:
		return structs.User{}, errors.New("invalid parameter")
	}

	// Convert the database row into a user struct
	user, err := ConvertRowToUser(query)
	if err != nil {
		return structs.User{}, errors.New("failed to convert")
	}

	if len(user) == 0 {
		return structs.User{}, errors.New("no users found")
	}

	return user[0], nil
}
