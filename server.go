package main

// this is my new Go plaground!! :)

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/brwhale/GoServer/controllers"
	"github.com/brwhale/GoServer/database"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	controllers.GenerateSecureCookie()
	database.Connect()
	defer database.Disconnect()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/content/{type}/{filename}", StaticHandler)
	router.HandleFunc("/", controllers.Index)
	router.HandleFunc("/login", controllers.Login)
	router.HandleFunc("/logout", controllers.Logout)
	router.HandleFunc("/register", controllers.Register)
	router.HandleFunc("/game", controllers.GameStart)
	router.HandleFunc("/game/{gamestate}", controllers.Game)
	router.HandleFunc("/commentlist", controllers.CommentList)
	router.HandleFunc("/boards", controllers.Boards)
	router.HandleFunc("/post-edit/{type}/{postid}", controllers.PostEdit)
	router.HandleFunc("/post-reply/{type}/{postid}", controllers.PostReply)
	router.HandleFunc("/test", controllers.Test)
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
