package main

// this is my new Go plaground!! :)

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/brwhale/GoServer/controllers"
	"github.com/brwhale/GoServer/database"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	database.Connect()
	defer database.Disconnect()
	database.GenerateSecureCookie()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/content/{type}/{filename}", StaticHandler)
	router.HandleFunc("/", controllers.Index)
	router.HandleFunc("/login", controllers.Login)
	router.HandleFunc("/logout", controllers.Logout)
	router.HandleFunc("/register", controllers.Register)
	router.HandleFunc("/game", controllers.GameStart)
	router.HandleFunc("/game/{gamestate}", controllers.Game)
	router.HandleFunc("/boards", controllers.Boards)
	router.HandleFunc("/user/{username}", controllers.User)
	router.HandleFunc("/post-edit/{type}/{postid}", controllers.PostEdit)
	router.HandleFunc("/post-reply/{type}/{postid}", controllers.PostReply)
	router.HandleFunc("/test", controllers.Test)

	log.Fatal(http.ListenAndServe(":8080", router))
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
