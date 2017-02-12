package controllers

import (
	"html/template"
	"net/http"

	"github.com/brwhale/GoServer/database"
	"github.com/brwhale/GoServer/drawing"
	"github.com/gorilla/mux"
)

// Page is some underused bs right now
type Page struct {
	CurrentUser string
	Title       string
}

// MessageBoard is a container for posts
type MessageBoard struct {
	CurrentUser string
	Posts       []database.Post
}

var layouts = template.Must(template.ParseGlob("html/layout/*"))
var templates = template.Must(layouts.ParseGlob("html/pages/*"))

// Index is the Home page
func Index(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "Index", Page{CurrentUser: GetSecureUsername(r)})
}

// GameStart is the landing page for the game
func GameStart(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "Start", Page{CurrentUser: GetSecureUsername(r)})
}

// Game is the dynamic game pages
func Game(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := Page{Title: vars["gamestate"], CurrentUser: GetSecureUsername(r)}
	templates.ExecuteTemplate(w, "Game", page)
}

// Test is a test of image generation
func Test(w http.ResponseWriter, r *http.Request) {
	drawing.GetImage(w)
}
