package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/brwhale/GoServer/util"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v2"
)

// database object
var db *sql.DB

// DbConfig is the database configuration format
type DbConfig struct {
	Username string
	Password string
	DbName   string
}

// Page is some underused bs right now
type Page struct {
	Title string
}

// Comment is the meat of it so far
type Comment struct {
	Author, Content                      string
	CreatedTime, EditedTime, UpdatedTime time.Time
	DisplayTime                          string
	Comments                             []Comment
	IsOwnComment                         bool
	ID, ParentID, PostID                 int
}

// Post is a container for comments
type Post struct {
	Author, Content                      string
	CreatedTime, EditedTime, UpdatedTime time.Time
	DisplayTime                          string
	Comments                             []Comment
	IsOwnPost                            bool
	ID                                   int
}

// MessageBoard is a container for posts
type MessageBoard struct {
	CurrentUser string
	Posts       []Post
}

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

// Connect to database
func Connect() {
	// read credentials from config file
	d := DbConfig{}
	b, err := ioutil.ReadFile("dbconfig.yaml")
	util.Check(err)
	err = yaml.Unmarshal(b, &d)
	util.Check(err)
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		d.Username, d.Password, d.DbName)
	// open the db
	db, err = sql.Open("postgres", dbinfo)
	util.Check(err)
}

// Disconnect the database
func Disconnect() {
	db.Close()
}

// Insert a Comment
func (comment Comment) Insert() {
	_, err := db.Exec("INSERT INTO comments(author,content,created,edited,updated,post_id,parent_comment) VALUES($1,$2,$3,$4,$5,$6,$7)", comment.Author, comment.Content, comment.CreatedTime, comment.EditedTime, comment.UpdatedTime, comment.PostID, comment.ParentID)
	util.Check(err)
	now := time.Now()
	_, err = db.Exec("UPDATE posts SET updated = $1 WHERE id = $2", now, comment.PostID)
	util.Check(err)
	if comment.ParentID > 0 {
		_, err = db.Exec("UPDATE comments SET updated = $1 WHERE id = $2", now, comment.ParentID)
		util.Check(err)
	}
}

// UpdateContent of a Comment
func (comment Comment) UpdateContent() {
	now := time.Now()
	_, err := db.Exec("UPDATE comments SET content = $1, updated = $2, edited = $2 WHERE id = $3", comment.Content, now, comment.ID)
	util.Check(err)
}

// GetComment gets a comment for verification
func GetComment(id int) Comment {
	row := db.QueryRow("SELECT author,content,created FROM comments WHERE id = $1", id)
	var comment Comment
	err := row.Scan(&comment.Author, &comment.Content, &comment.CreatedTime)
	util.Check(err)
	return comment
}

// GetComments gets the comments
func GetComments() []Comment {
	rows, err := db.Query("SELECT id,author,content,created,edited,updated,post_id,parent_comment FROM comments ORDER BY updated DESC")
	util.Check(err)
	// reform rows into comments
	var Comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.Author, &comment.Content, &comment.CreatedTime, &comment.EditedTime, &comment.UpdatedTime, &comment.PostID, &comment.ParentID)
		util.Check(err)
		comment.DisplayTime = friendlyString(time.Since(comment.CreatedTime))
		Comments = append(Comments, comment)
	}
	err = rows.Err()
	util.Check(err)

	return Comments
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

func friendlyString(duration time.Duration) string {
	if duration.Hours() >= 48 {
		return fmt.Sprintf("%.0f days ago", duration.Hours()/24)
	}
	if duration.Hours() >= 24 {
		return "1 day ago"
	}
	if duration.Hours() >= 2 {
		return fmt.Sprintf("%.0f hours ago", duration.Hours())
	}
	if duration.Hours() >= 1 {
		return "1 hour ago"
	}
	if duration.Minutes() >= 2 {
		return fmt.Sprintf("%.0f minutes ago", duration.Minutes())
	}
	if duration.Minutes() >= 1 {
		return "1 minute ago"
	}
	return "a couple seconds ago"
}
