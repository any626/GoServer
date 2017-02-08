package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/brwhale/GoServer/database"
	"github.com/brwhale/GoServer/drawing"
	"github.com/brwhale/GoServer/util"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
)

var sc *securecookie.SecureCookie
var layouts = template.Must(template.ParseGlob("html/layout/*"))
var templates = template.Must(layouts.ParseGlob("html/pages/*"))

// GenerateSecureCookie gets new crypto
func GenerateSecureCookie() {
	var hashKey = securecookie.GenerateRandomKey(64)
	var blockKey = securecookie.GenerateRandomKey(32)
	sc = securecookie.New(hashKey, blockKey)
}

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
		id := database.ValidatePassword(Username, Password)
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
		util.Check(err)
		id := database.User{Name: Username, Hash: string(hash)}.Insert()
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
	page := &database.Page{Title: vars["gamestate"]}
	templates.ExecuteTemplate(w, "Game", page)
}

// CommentList is the wip bulletin board
func CommentList(w http.ResponseWriter, r *http.Request) {
	Author := GetSecureUsername(r)
	Content := r.FormValue("Content")
	if Author != "" && Content != "" {
		// add new comment
		now := time.Now()
		database.Comment{
			Author:      Author,
			Content:     Content,
			CreatedTime: now,
			EditedTime:  now,
		}.Insert()
	}
	var thePost database.Post
	thePost.Comments = database.GetComments()
	templates.ExecuteTemplate(w, "CommentList", thePost)
}

func findRoot(data *map[int]*database.Comment, id int) int {
	if elem, ok := (*data)[id]; ok {
		if (*elem).ParentID != 0 {
			findRoot(data, elem.ParentID)
		} else {
			return id
		}
	}
	return 0
}

func recGetComs(list []database.Comment, parentID int) []database.Comment {
	var comments []database.Comment
	for _, c := range list {
		if c.ParentID == parentID {
			c.Comments = recGetComs(list, c.ID)
			comments = append(comments, c)
		}
	}
	return comments
}

// Boards is the wip bulletin board
func Boards(w http.ResponseWriter, r *http.Request) {
	Author := GetSecureUsername(r)
	Content := r.FormValue("Content")
	if Author != "" && Content != "" {
		// add new post
		now := time.Now()
		database.Post{
			Author:      Author,
			Content:     Content,
			CreatedTime: now,
			EditedTime:  now,
			UpdatedTime: now,
		}.Insert()
	}
	var mainpage database.MessageBoard
	mainpage.CurrentUser = Author
	mainpage.Posts = database.GetPosts()
	for index := range mainpage.Posts {
		mainpage.Posts[index].IsOwnPost = Author == mainpage.Posts[index].Author
	}
	comments := database.GetComments()
	for index := range comments {
		comments[index].IsOwnComment = Author == comments[index].Author
	}
	for _, comment := range comments {
		// sorted already
		if comment.ParentID == 0 {
			for pindex := range mainpage.Posts {
				if mainpage.Posts[pindex].ID == comment.PostID {
					comment.Comments = recGetComs(comments, comment.ID)
					mainpage.Posts[pindex].Comments = append(mainpage.Posts[pindex].Comments, comment)
					break
				}
			}
		}
	}

	util.PrettyPrint(mainpage)
	templates.ExecuteTemplate(w, "MessageBoard", mainpage)
}

// PostEdit edits posts
func PostEdit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, _ := strconv.Atoi(vars["postid"])
	thingType := vars["type"]
	Author := GetSecureUsername(r)
	Content := r.FormValue("Content")
	if Author != "" && Content != "" {
		if thingType == "comment" {
			oldcomment := database.GetComment(postID)
			if oldcomment.Author == Author {
				database.Comment{
					ID:      postID,
					Content: Content,
				}.UpdateContent()
			}
		} else if thingType == "post" {
			oldpost := database.GetPost(postID)
			if oldpost.Author == Author {
				database.Post{
					ID:      postID,
					Content: Content,
				}.UpdateContent()
			}
		}
	}
	http.Redirect(w, r, "/boards", 302)
}

// PostReply replies to posts
func PostReply(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, _ := strconv.Atoi(vars["postid"])
	thingType := vars["type"]
	Author := GetSecureUsername(r)
	Content := r.FormValue("Content")
	if thingType == "post" {
		if Author != "" && Content != "" {
			now := time.Now()
			database.Comment{
				Author:      Author,
				Content:     Content,
				CreatedTime: now,
				EditedTime:  now,
				UpdatedTime: now,
				PostID:      postID,
			}.Insert()
		}
	} else {
		if Author != "" && Content != "" {
			parent := postID
			postID, _ = strconv.Atoi(thingType)
			now := time.Now()
			database.Comment{
				Author:      Author,
				Content:     Content,
				CreatedTime: now,
				EditedTime:  now,
				UpdatedTime: now,
				PostID:      postID,
				ParentID:    parent,
			}.Insert()
		}
	}
	http.Redirect(w, r, "/boards", 302)
}

// Test is a test of image generation
func Test(w http.ResponseWriter, r *http.Request) {
	drawing.GetImage(w)
}
