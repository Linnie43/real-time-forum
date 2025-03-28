package structs

import "time"

type Login struct {
	Entry    string `json:"username"`
	Password string `json:"password"`
}

type Comment struct {
	Id      int       `json:"id"`
	Post_id int       `json:"post_id"`
	User_id int       `json:"user_id"`
	Content string    `json:"content"`
	Date    time.Time `json:"date"`
}

type Message struct {
	Id          int    `json:"id"`
	Sender_id   int    `json:"sender_id"`
	Receiver_id int    `json:"receiver_id"`
	Content     string `json:"content"`
	Date        string `json:"date"`
	Msg_type    string `json:"msg_type"`
}

type Post struct {
	Id       int       `json:"id"`
	User_id  int       `json:"user_id"`
	Category string    `json:"category"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Date     time.Time `json:"date"`
}

type User struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Gender    string `json:"gender"`
	DOB       string `json:"dob"`
	Password  string `json:"password"`
}

type Chat struct {
	User_one int
	User_two int
	Time     int
}

type OnlineUsers struct {
	UserIds  []int  `json:"user_ids"`
	Msg_type string `json:"msg_type"`
}
