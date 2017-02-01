package main
// this is a basic web service in Go! :)

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"io"
	"time"
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
	router.HandleFunc("/css/{path}", CssHandler)
	router.HandleFunc("/", Index)
	router.HandleFunc("/game", GameStart)
	router.HandleFunc("/game/{gamestate}", Game)
	router.HandleFunc("/CommentList", CommentList)
	log.Fatal (http.ListenAndServe(":8080", router))
}

func CssHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	static_file := vars["path"]
	if len(static_file) != 0 {
		f, err := http.Dir("css/").Open(static_file)
		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, static_file, time.Now(), content)
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
