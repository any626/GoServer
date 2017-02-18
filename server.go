package main

// this is my new Go plaground!

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
	var db database.KataDB
	db.Connect()
	defer db.Disconnect()
	db.GenerateSecureCookie()
	kc := controllers.NewController(&db)
	// set routing
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/content/{type}/{filename}", StaticHandler)
	router.HandleFunc("/", kc.Index)
	router.HandleFunc("/login", kc.Login)
	router.HandleFunc("/logout", kc.Logout)
	router.HandleFunc("/register", kc.Register)
	router.HandleFunc("/game", kc.GameStart)
	router.HandleFunc("/game/{gamestate}", kc.Game)
	router.HandleFunc("/boards", kc.Boards)
	router.HandleFunc("/user/{username}", kc.User)
	router.HandleFunc("/post-edit/{type}/{postid}", kc.PostEdit)
	router.HandleFunc("/post-reply/{type}/{postid}", kc.PostReply)
	router.HandleFunc("/test", kc.Test)
	// redirect to https
	unsecureserver := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Connection", "close")
			url := "https://" + req.Host + req.URL.String()
			http.Redirect(w, req, url, http.StatusMovedPermanently)
		}),
	}
	go func() { log.Fatal(unsecureserver.ListenAndServe()) }()
	// run https server
	secureserver := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}
	log.Fatal(secureserver.ListenAndServeTLS("server.crt", "server.key"))
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
