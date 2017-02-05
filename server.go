package main

// this is a basic web service in Go! :)

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var layouts = template.Must(template.ParseGlob("html/layout/*"))
var templates = template.Must(layouts.ParseGlob("html/pages/*"))

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	Connect()
	defer Disconnect()
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
