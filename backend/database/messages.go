package database

import (
	"database/sql"
	"errors"
	"strconv"

	"real-time-forum/backend/database/structs"
)

// Attempts to insert a new message into the database
func NewMessage(path string, message structs.Message) error {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return err
	}

	defer db.Close()

	//Executes the insert statement
	_, err = db.Exec(AddMessage, message.Sender_id, message.Receiver_id, message.Content, message.Date)
	if err != nil {
		return err
	}

	err = UpdateChatTime(message.Sender_id, message.Receiver_id, db)
	if err != nil {
		return err
	}

	return nil
}

// Converts message table query results into an array of message structs
func ConvertRowToMessage(rows *sql.Rows) ([]structs.Message, error) {
	var messages []structs.Message

	//Loops through the rows provided
	for rows.Next() {
		var m structs.Message

		//Stores the row data in a temporary message struct
		err := rows.Scan(&m.Id, &m.Sender_id, &m.Receiver_id, &m.Content, &m.Date)
		if err != nil {
			break
		}

		//Appends the temporary struct to the array
		messages = append(messages, m)
	}

	//Returns an error if no rows are provided
	// if len(messages) == 0 {
	// 	return []models.Message{}, errors.New("no row provided")
	// }

	return messages, nil
}

// Finds chat messages between users
func FindChatMessages(path, sender, receiver string) ([]structs.Message, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return []structs.Message{}, errors.New("failed to open database")
	}

	defer db.Close()

	//Converts sender and receiver ids to integers
	s, err := strconv.Atoi(sender)
	if err != nil {
		return []structs.Message{}, errors.New("sender id must be an integer")
	}

	r, err := strconv.Atoi(receiver)
	if err != nil {
		return []structs.Message{}, errors.New("receiver id must be an integer")
	}

	//Searches database for all messages between the two users
	q, err := db.Query(GetAllChatMessage, s, r, r, s)
	if err != nil {
		return []structs.Message{}, errors.New("could not find chat messages")
	}

	//Converts rows to an array of message structs
	messages, err := ConvertRowToMessage(q)
	if err != nil {
		return []structs.Message{}, errors.New("failed to convert")
	}

	return messages, nil
}
