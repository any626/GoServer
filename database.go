package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

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
	Author, Content string
	CreatedTime     time.Time
}

// Post is a container for comments
type Post struct {
	Comments []Comment
}

// Connect to database
func Connect() {
	// read credentials from config file
	d := DbConfig{}
	b, err := ioutil.ReadFile("dbconfig.yaml")
	checkErr(err)
	err = yaml.Unmarshal(b, &d)
	checkErr(err)
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		d.Username, d.Password, d.DbName)
	// open the db
	db, err = sql.Open("postgres", dbinfo)
	checkErr(err)
}

// Disconnect the database
func Disconnect() {
	db.Close()
}

// Insert a Comment
func (comment Comment) Insert() {
	_, err := db.Exec("INSERT INTO comments(author,content,created) VALUES($1,$2,$3)", comment.Author, comment.Content, comment.CreatedTime)
	checkErr(err)
}

// GetComments gets the comments
func GetComments() []Comment {
	rows, err := db.Query("SELECT author,content,created FROM comments")
	checkErr(err)
	// reform rows into comments
	var Comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.Author, &comment.Content, &comment.CreatedTime)
		checkErr(err)
		Comments = append(Comments, comment)
	}
	err = rows.Err()
	checkErr(err)

	return Comments
}
