package controllers

import (
	"html/template"
	"net/http"

	"github.com/brwhale/GoServer/database"
	"github.com/brwhale/GoServer/drawing"
	"github.com/gorilla/mux"
)

type validation struct {
	Errors bool
}

var layouts = template.Must(template.ParseGlob("html/layout/*"))
var templates = template.Must(layouts.ParseGlob("html/pages/*"))

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
	page := &database.Page{Title: vars["gamestate"]}
	templates.ExecuteTemplate(w, "Game", page)
}

// Test is a test of image generation
func Test(w http.ResponseWriter, r *http.Request) {
	drawing.GetImage(w)
}
