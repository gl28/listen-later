package models

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials error = errors.New("Invalid login credentials")
var ErrUserAlreadyExists error = errors.New("A user with that email already exists")

type User struct {
	id int
	email string
	key string
}

func generateKey() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func RegisterNewUser(email, password string) error {
	var existing_id int
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&existing_id)
	if err == nil {
		return ErrUserAlreadyExists
	} else if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO users (email, key, hash) VALUES ($1, $2, $3)")
	if err != nil {
		log.Fatal(err)
	}
	key := generateKey()
	_, err = stmt.Exec(email, key, hash)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}