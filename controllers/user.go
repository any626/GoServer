package controllers

import (
	"net/http"

	"github.com/brwhale/GoServer/database"
	"github.com/gorilla/mux"
)

type userView struct {
	User        database.User
	CurrentUser string
}

// User edits posts
func User(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	currentUser := GetSecureUsername(r)
	user := database.User{Name: vars["username"]}
	user.Comments = user.GetComments()
	for index := range user.Comments {
		user.Comments[index].IsOwnComment = currentUser == user.Comments[index].Author
	}
	templates.ExecuteTemplate(w, "User", userView{User: user, CurrentUser: currentUser})
}
