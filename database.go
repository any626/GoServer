package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type DbConfig struct {
	Username string
	Password string
	DbName   string
}

var db *sql.DB

type Page struct {
	Title string
}
type Comment struct {
	Author, Content string
	CreatedTime     time.Time
}
type Post struct {
	Comments []Comment
}

func Connect() {
	d := DbConfig{}
	b, err := ioutil.ReadFile("dbconfig.yaml")
	checkErr(err)
	err = yaml.Unmarshal(b, &d)
	checkErr(err)
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		d.Username, d.Password, d.DbName)
	db, err = sql.Open("postgres", dbinfo)
	checkErr(err)
}

func Disconnect() {
	db.Close()
}

func (comment Comment) Insert() {
	_, err := db.Exec("INSERT INTO comments(author,content,created) VALUES($1,$2,$3)", comment.Author, comment.Content, comment.CreatedTime)
	checkErr(err)
}

func GetComments() []Comment {
	rows, err := db.Query("SELECT author,content,created FROM comments")
	checkErr(err)
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
