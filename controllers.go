package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "Index", nil)
}

func GameStart(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "Start", nil)
}

func Game(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := &Page{Title: vars["gamestate"]}
	templates.ExecuteTemplate(w, "Game", page)
}

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

func Test(w http.ResponseWriter, r *http.Request) {
	GetImage(w)
}
