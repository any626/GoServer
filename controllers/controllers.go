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
	Posts       []*database.Post
}

// KataController is the main struct for the pages
type KataController struct {
	layouts, templates *template.Template
	DB                 *database.KataDB
}

// NewController makes a new controller, inject a db into it
func NewController(db *database.KataDB) *KataController {
	var theController KataController
	theController.DB = db
	theController.layouts = template.Must(template.ParseGlob("html/layout/*"))
	theController.templates = template.Must(theController.layouts.ParseGlob("html/pages/*"))
	return &theController
}

// Index is the Home page
func (c *KataController) Index(w http.ResponseWriter, r *http.Request) {
	c.templates.ExecuteTemplate(w, "Index", Page{CurrentUser: c.GetSecureUsername(r)})
}

// GameStart is the landing page for the game
func (c *KataController) GameStart(w http.ResponseWriter, r *http.Request) {
	c.templates.ExecuteTemplate(w, "Start", Page{CurrentUser: c.GetSecureUsername(r)})
}

// Game is the dynamic game pages
func (c *KataController) Game(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := Page{Title: vars["gamestate"], CurrentUser: c.GetSecureUsername(r)}
	c.templates.ExecuteTemplate(w, "Game", page)
}

// Test is a test of image generation
func (c *KataController) Test(w http.ResponseWriter, r *http.Request) {
	drawing.GetImage(w)
}
