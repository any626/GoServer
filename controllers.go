package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Index is the Home page
func Index(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "Index", nil)
}

// GameStart is the landing page for the game
func GameStart(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "Start", nil)
}

// Game is the dynamic game pages
func Game(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := &Page{Title: vars["gamestate"]}
	templates.ExecuteTemplate(w, "Game", page)
}

// CommentList is the wip bulletin board
func CommentList(w http.ResponseWriter, r *http.Request) {
	Author := r.FormValue("Author")
	Content := r.FormValue("Content")
	if Author != "" && Content != "" {
		Comment{
			Author:      Author,
			Content:     Content,
			CreatedTime: time.Now(),
		}.Insert()
	}
	var thePost Post
	thePost.Comments = GetComments()
	templates.ExecuteTemplate(w, "CommentList", thePost)
}

// Test is a test of image generation
func Test(w http.ResponseWriter, r *http.Request) {
	GetImage(w)
}
