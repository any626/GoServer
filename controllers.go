package main

import (
	"net/http"
	"time"

	"fmt"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Index is the Home page
func Index(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "Index", nil)
}

// FinalizeLogin generates cookie and save session info
func FinalizeLogin(username string, w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:  "login-name",
		Value: username,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", 302)
}

type validation struct {
	Errors bool
}

// Login is the login page for the site
func Login(w http.ResponseWriter, r *http.Request) {
	Username := r.FormValue("Username")
	Password := r.FormValue("Password")
	val := validation{Errors: false}
	if Username != "" && Password != "" {
		id := ValidatePassword(Username, Password)
		if id >= 0 {
			FinalizeLogin(Username, w, r)
		}
		val.Errors = true
	}
	templates.ExecuteTemplate(w, "Login", val)
}

// Register is the register page for the site
func Register(w http.ResponseWriter, r *http.Request) {
	Username := r.FormValue("Username")
	Password := r.FormValue("Password")
	val := validation{Errors: false}
	if Username != "" && Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
		checkErr(err)
		id := User{Name: Username, Hash: string(hash)}.Insert()
		if id >= 0 {
			FinalizeLogin(Username, w, r)
		}
		fmt.Println(id)
		val.Errors = true
	}
	templates.ExecuteTemplate(w, "Register", val)
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
	Author := ""
	Content := r.FormValue("Content")
	if cookie, err := r.Cookie("login-name"); err == nil {
		Author = cookie.Value
	}
	if Author != "" && Content != "" {
		// add new comment
		Comment{
			Author:      Author,
			Content:     Content,
			CreatedTime: time.Now(),
		}.Insert()
	}
	var thePost Post
	thePost.CurrentUser = Author
	thePost.Comments = GetComments()
	templates.ExecuteTemplate(w, "CommentList", thePost)
}

// Test is a test of image generation
func Test(w http.ResponseWriter, r *http.Request) {
	GetImage(w)
}
