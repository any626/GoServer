package main
// this is a basic web service in Go! :)

import (
	"html/template"
	"log"
	"net/http"
	"io"
	"time"

	"github.com/gorilla/mux"
)

type Page struct {
	Title string
}
type Comment struct {
	Author string
	Content string
	CreatedTime time.Time
}
type Post struct {
	Comments []Comment
}

var layouts = template.Must(template.ParseGlob("html/layout/*"))
var templates = template.Must(layouts.ParseGlob("html/pages/*"))

var thePost Post

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/content/{type}/{filename}", StaticHandler)
	router.HandleFunc("/", Index)
	router.HandleFunc("/game", GameStart)
	router.HandleFunc("/game/{gamestate}", Game)
	router.HandleFunc("/CommentList", CommentList)
	log.Fatal (http.ListenAndServe(":8080", router))
}

func StaticHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	if len(vars["filename"]) != 0 {
		f, err := http.Dir("content/"+vars["type"]+"/").Open(vars["filename"])
		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, vars["filename"], time.Now(), content)
			return
		}
	}
	http.NotFound(w, req)
}

func Index(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "Index", nil)
}

func GameStart(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "Start", nil)
}

func Game(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := &Page{Title: vars["gamestate"]}
	templates.ExecuteTemplate(w, "test", page)
}

func CommentList(w http.ResponseWriter, r *http.Request) {
	Author := r.FormValue("Author")
	Content := r.FormValue("Content")
	if Author != "" && Content != "" {
		comment := Comment{Author: Author, Content: Content, CreatedTime: time.Now()}
		thePost.Comments = append(thePost.Comments, comment)
	}
	templates.ExecuteTemplate(w, "CommentList", thePost)
}
