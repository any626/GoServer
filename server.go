package main

// this is my Go plaground!! :)

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	_ "github.com/lib/pq"
)

var layouts = template.Must(template.ParseGlob("html/layout/*"))
var templates = template.Must(layouts.ParseGlob("html/pages/*"))
var sc *securecookie.SecureCookie

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var hashKey = securecookie.GenerateRandomKey(64)
	var blockKey = securecookie.GenerateRandomKey(32)
	sc = securecookie.New(hashKey, blockKey)
	Connect()
	defer Disconnect()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/content/{type}/{filename}", StaticHandler)
	router.HandleFunc("/", Index)
	router.HandleFunc("/login", Login)
	router.HandleFunc("/logout", Logout)
	router.HandleFunc("/register", Register)
	router.HandleFunc("/game", GameStart)
	router.HandleFunc("/game/{gamestate}", Game)
	router.HandleFunc("/commentlist", CommentList)
	router.HandleFunc("/test", Test)
	log.Fatal(http.ListenAndServe(":8080", LowerCaseURI(router)))
}

// LowerCaseURI lowercases all the urls so it seems case insensitive
func LowerCaseURI(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.ToLower(r.URL.Path)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// StaticHandler handles static content such as images, css, javascript, etc
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
