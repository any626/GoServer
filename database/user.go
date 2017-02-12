package database

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User is a user of the site!!
type User struct {
	ID                int
	Name, Email, Hash string
	CreatedTime       time.Time
	Banned            bool
}

// Insert a User
func (user User) Insert() int {
	var id int
	err := db.QueryRow("INSERT INTO users(name,email,hash,created) VALUES($1,$2,$3,$4) RETURNING id", user.Name, user.Email, user.Hash, user.CreatedTime).Scan(&id)
	if err != nil {
		return -1
	}
	return int(id)
}

// ValidatePassword validates the pw!
func ValidatePassword(username, password string) int {
	row := db.QueryRow("SELECT id,hash FROM users WHERE name = $1", username)
	var user User
	err := row.Scan(&user.ID, &user.Hash)
	if err == nil {
		err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
		if err == nil {
			return user.ID
		}
	}
	return -1
}
