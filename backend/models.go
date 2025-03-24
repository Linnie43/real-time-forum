package main

import (
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        string
	Nickname  string
	Email     string
	Password  string
	Age       int
	Gender    string
	FirstName string
	LastName  string
}

// Post represents a forum post
type Post struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
}

// HashPassword securely hashes the user's password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash verifies a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// NewUser creates a new user with a hashed password
func NewUser(nickname, email, password string, age int, gender, firstName, lastName string) (*User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	id, _ := uuid.NewV4()

	return &User{
		ID:        id.String(),
		Nickname:  nickname,
		Email:     email,
		Password:  hashedPassword,
		Age:       age,
		Gender:    gender,
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}
