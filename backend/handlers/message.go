package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"real-time-forum/backend/database"
	"real-time-forum/backend/database/structs"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	//Prevents the endpoint being called by other url paths
	if r.URL.Path != "/message" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	//Checks whether it is a POST or GET request
	switch r.Method {
	case "GET":
		cookie, err := r.Cookie("session")
		if err != nil {
			return
		}

		foundVal := cookie.Value

		user, err := database.CurrentUser("database.db", foundVal)
		if err != nil {
			return
		}

		sender := strconv.Itoa(user.Id)

		//Grabs the receiver id from the url
		receiver := r.URL.Query().Get("receiver")

		//Grabs the offset from the url
		offset := r.URL.Query().Get("offset")
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		//Makes sure neither are empty
		if receiver == "" || offset == "" {
			http.Error(w, "400 bad request", http.StatusBadRequest)
			return
		}

		//Gets the messages from the database
		messages, err := database.FindChatMessages("database.db", sender, receiver)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		// Keep only the 10 messages after the offset starting from the end of the array
		if len(messages) > offsetInt+10 { // If there are more than 10 messages after the offset
			messages = messages[len(messages)-10-offsetInt : len(messages)-offsetInt]
		} else if len(messages) > offsetInt { // If there are less than 10 messages after the offset
			messages = messages[0 : len(messages)-offsetInt]
		} else {
			messages = []structs.Message{}
		}

		//Marshals the array of message structs to a json object
		resp, err := json.Marshal(messages)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}

		//Writes the json object to the frontend
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	case "POST":
		var newMessage structs.Message

		//Decodes the request body into the message struct
		//Returns a bad request if there's an error
		err := json.NewDecoder(r.Body).Decode(&newMessage)
		if err != nil {
			http.Error(w, "400 bad request.", http.StatusBadRequest)
			return
		}

		//Attemps to add the new message to the database
		err = database.NewMessage("database.db", newMessage)
		if err != nil {
			http.Error(w, "500 internal server error", http.StatusInternalServerError)
			return
		}
	default:
		//Prevents the use of other request types
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
