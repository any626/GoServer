package main

// this is a basic web service in Go! :)

import (
	"html/template"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Page struct {
	Title string
}
type Comment struct {
	Author, Content string
	CreatedTime     time.Time
}
type Post struct {
	Comments []Comment
}

type Circle struct {
	X, Y, R float64
}

func (c *Circle) Brightness(x, y float64) uint8 {
	var dx, dy float64 = c.X - x, c.Y - y
	d := math.Sqrt(dx*dx+dy*dy) / c.R
	if d > 1 {
		return 0
	} else {
		return 255
	}
}

func GetImage(writer http.ResponseWriter) {
	var w, h int = 280, 240
	var hw, hh float64 = float64(w / 2), float64(h / 2)
	r := 40.0
	θ := 2 * math.Pi / 3
	cr := &Circle{hw - r*math.Sin(0), hh - r*math.Cos(0), 60}
	cg := &Circle{hw - r*math.Sin(θ), hh - r*math.Cos(θ), 60}
	cb := &Circle{hw - r*math.Sin(-θ), hh - r*math.Cos(-θ), 60}

	m := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			c := color.RGBA{
				cr.Brightness(float64(x), float64(y)),
				cg.Brightness(float64(x), float64(y)),
				cb.Brightness(float64(x), float64(y)),
				255,
			}
			m.Set(x, y, c)
		}
	}
	png.Encode(writer, m)
}

func Test(w http.ResponseWriter, r *http.Request) {
	GetImage(w)
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
	router.HandleFunc("/test", Test)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func StaticHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	if len(vars["filename"]) != 0 {
		f, err := http.Dir("content/" + vars["type"] + "/").Open(vars["filename"])
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
	templates.ExecuteTemplate(w, "Game", page)
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
