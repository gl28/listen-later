package models

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials error = errors.New("Invalid login credentials")
var ErrUserAlreadyExists error = errors.New("A user with that email already exists")

type User struct {
	Id int
	Email string
	Key string
	ArticleCount int
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
		return err
	}
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("INSERT INTO users (email, key, hash) VALUES ($1, $2, $3)")
	if err != nil {
		return err
	}
	key := generateKey()
	_, err = stmt.Exec(email, key, hash)
	if err != nil {
		return err
	}
	return nil
}

func AuthenticateUser(email, password string) (*User, error) {
	var (
		id int
		key string
		hash string
		articleCount int
	)
	err := db.QueryRow("SELECT id, key, hash, article_count FROM users WHERE email = $1", email).Scan(&id, &key, &hash, &articleCount)
	if err == sql.ErrNoRows {
		return nil, ErrInvalidCredentials
	} else if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	return &User{Id: id, Email: email, Key: key, ArticleCount: articleCount}, nil
}

func GetUserById(userId int) (*User, error) {
	var (
		key string
		email string
		articleCount int
	)
	err := db.QueryRow("SELECT email, key, article_count FROM users WHERE id = $1", userId).Scan(&email, &key, &articleCount)
	if err != nil {
		return nil, err
	}
	return &User{Id: userId, Email: email, Key: key, ArticleCount: articleCount}, nil
}


/* THIS FUNCTION NOT BEING USED BECAUSE AWS SETS THE AUDIO URL AUTOMATICALLY NOW
func IncrementArticleCount(userId int) error {
	stmt, err := db.Prepare("UPDATE users SET article_count = article_count + 1 WHERE id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userId)
	return err
}
*/