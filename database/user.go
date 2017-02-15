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
	Comments          []*Comment
}

// Insert a User
func (user User) Insert() error {
	_, err := db.Exec("INSERT INTO users(name,email,hash,created) VALUES($1,$2,$3,$4) RETURNING id", user.Name, user.Email, user.Hash, user.CreatedTime)
	return err
}

// ValidatePassword validates the pw!
func ValidatePassword(username, password string) bool {
	row := db.QueryRow("SELECT id,hash FROM users WHERE name = $1", username)
	var user User
	err := row.Scan(&user.ID, &user.Hash)
	if err == nil {
		err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
		if err == nil {
			return true
		}
	}
	return false
}

// GetComments gets the comments
func (user *User) GetComments() ([]*Comment, error) {
	var Comments []*Comment
	rows, err := db.Query("SELECT id,author,content,created,edited,updated,post_id,parent_comment FROM comments WHERE author = $1 ORDER BY updated DESC", user.Name)
	if err != nil {
		return Comments, err
	}
	// reform rows into comments
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.Author, &comment.Content, &comment.CreatedTime, &comment.EditedTime, &comment.UpdatedTime, &comment.PostID, &comment.ParentID)
		if err != nil {
			return Comments, err
		}
		comment.DisplayTime = friendlyString(time.Since(comment.CreatedTime))
		Comments = append(Comments, &comment)
	}
	err = rows.Err()

	return Comments, err
}
