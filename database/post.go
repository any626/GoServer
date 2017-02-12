package database

import (
	"time"

	"github.com/brwhale/GoServer/util"
)

// Post is a container for comments
type Post struct {
	Author, Content                      string
	CreatedTime, EditedTime, UpdatedTime time.Time
	DisplayTime                          string
	Comments                             []Comment
	IsOwnPost                            bool
	ID                                   int
}

// Insert a Post
func (post Post) Insert() {
	_, err := db.Exec("INSERT INTO posts(author,content,created,edited,updated) VALUES($1,$2,$3,$4,$5)", post.Author, post.Content, post.CreatedTime, post.EditedTime, post.UpdatedTime)
	util.Check(err)
}

// UpdateContent of a Post
func (post Post) UpdateContent() {
	now := time.Now()
	_, err := db.Exec("UPDATE posts SET content = $1, updated = $2, edited = $2 WHERE id = $3", post.Content, now, post.ID)
	util.Check(err)
}

// GetPosts gets the posts
func GetPosts() []Post {
	rows, err := db.Query("SELECT id,author,content,created,updated FROM posts ORDER BY updated DESC")
	util.Check(err)
	// reform rows into posts
	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Author, &post.Content, &post.CreatedTime, &post.UpdatedTime)
		util.Check(err)
		post.DisplayTime = friendlyString(time.Since(post.CreatedTime))
		posts = append(posts, post)
	}
	err = rows.Err()
	util.Check(err)

	return posts
}

// GetPost gets a post
func GetPost(id int) Post {
	row := db.QueryRow("SELECT author,content,created FROM posts WHERE id = $1", id)
	var post Post
	err := row.Scan(&post.Author, &post.Content, &post.CreatedTime)
	util.Check(err)
	return post
}
