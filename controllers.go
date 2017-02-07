package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// GetSecureUsername from secure cookie
func GetSecureUsername(r *http.Request) string {
	Username := ""
	if cookie, err := r.Cookie("login-name"); err == nil {
		value := make(map[string]string)
		if err = sc.Decode("login-name", cookie.Value, &value); err == nil {
			Username = value["username"]
		}
	}
	return Username
}

// FinalizeLogin generates cookie and save session info
func FinalizeLogin(username string, w http.ResponseWriter, r *http.Request) {
	if username == "" {
		cookie := &http.Cookie{
			Name:    "login-name",
			Value:   "",
			Path:    "/",
			Expires: time.Now(),
		}
		http.SetCookie(w, cookie)
	} else {
		value := map[string]string{
			"username": username,
		}
		if encoded, err := sc.Encode("login-name", value); err == nil {
			cookie := &http.Cookie{
				Name:  "login-name",
				Value: encoded,
				Path:  "/",
			}
			http.SetCookie(w, cookie)
		}
	}
	http.Redirect(w, r, "/", 302)
}

type validation struct {
	Errors bool
}

// Index is the Home page
func Index(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "Index", nil)
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

// Logout is the logout page for the site
func Logout(w http.ResponseWriter, r *http.Request) {
	FinalizeLogin("", w, r)
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
	Author := GetSecureUsername(r)
	Content := r.FormValue("Content")
	if Author != "" && Content != "" {
		// add new comment
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

// Boards is the wip bulletin board
func Boards(w http.ResponseWriter, r *http.Request) {
	Author := GetSecureUsername(r)
	Content := r.FormValue("Content")
	if Author != "" && Content != "" {
		// add new comment
		Post{
			Author:      Author,
			Content:     Content,
			CreatedTime: time.Now(),
		}.Insert()
	}
	var mainpage MessageBoard
	mainpage.CurrentUser = Author
	mainpage.Posts = GetPosts(Author)
	templates.ExecuteTemplate(w, "MessageBoard", mainpage)
}

// PostEdit edits posts
func PostEdit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, _ := strconv.Atoi(vars["postid"])
	Author := GetSecureUsername(r)
	Content := r.FormValue("Content")
	if Author != "" && Content != "" {
		oldpost := GetPost(postID)
		if oldpost.Author == Author {
			// update comment
			Post{
				ID:      postID,
				Content: Content,
			}.UpdateContent()
		}
	}
	http.Redirect(w, r, "/boards", 302)
}

// PostReply replies to posts
func PostReply(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, _ := strconv.Atoi(vars["postid"])
	Author := GetSecureUsername(r)
	Content := r.FormValue("Content")
	if Author != "" && Content != "" {
		Comment{
			Author:      Author,
			Content:     Content,
			CreatedTime: time.Now(),
			PostID:      postID,
		}.Insert()

	}
	http.Redirect(w, r, "/boards", 302)
}

// Test is a test of image generation
func Test(w http.ResponseWriter, r *http.Request) {
	GetImage(w)
}
