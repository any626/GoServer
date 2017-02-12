package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/brwhale/GoServer/database"
	"github.com/brwhale/GoServer/util"
	"golang.org/x/crypto/bcrypt"
)

type validation struct {
	CurrentUser string
	Errors      bool
}

// Login is the login page for the site
func Login(w http.ResponseWriter, r *http.Request) {
	Username := r.FormValue("Username")
	Password := r.FormValue("Password")
	val := validation{Errors: false, CurrentUser: GetSecureUsername(r)}
	if Username != "" && Password != "" {
		id := database.ValidatePassword(Username, Password)
		if id >= 0 {
			FinalizeLogin(Username, w, r)
		}
		val.Errors = true
	}
	templates.ExecuteTemplate(w, "Login", val)
}

// Logout is the logout page for the site
func Logout(w http.ResponseWriter, r *http.Request) {
	FinalizeLogin("", w, r)
}

// Register is the register page for the site
func Register(w http.ResponseWriter, r *http.Request) {
	Username := r.FormValue("Username")
	Password := r.FormValue("Password")
	val := validation{Errors: false, CurrentUser: GetSecureUsername(r)}
	if Username != "" && Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
		util.Check(err)
		id := database.User{Name: Username, Hash: string(hash)}.Insert()
		if id >= 0 {
			FinalizeLogin(Username, w, r)
		}
		fmt.Println(id)
		val.Errors = true
	}
	templates.ExecuteTemplate(w, "Register", val)
}

// FinalizeLogin generates cookie and save session info
func FinalizeLogin(username string, w http.ResponseWriter, r *http.Request) {
	if username == "" {
		cookie := &http.Cookie{
			Name:    "login-name",
			Value:   "",
			Path:    "/",
			Expires: time.Now(),
		}
		http.SetCookie(w, cookie)
	} else {
		value := map[string]string{
			"username": username,
		}
		if encoded, err := database.SecureEncode("login-name", value); err == nil {
			cookie := &http.Cookie{
				Name:  "login-name",
				Value: encoded,
				Path:  "/",
			}
			http.SetCookie(w, cookie)
		}
	}
	http.Redirect(w, r, "/", 302)
}

// GetSecureUsername from secure cookie
func GetSecureUsername(r *http.Request) string {
	Username := ""
	if cookie, err := r.Cookie("login-name"); err == nil {
		value := make(map[string]string)
		if err = database.SecureDecode("login-name", cookie.Value, &value); err == nil {
			Username = value["username"]
		}
	}
	return Username
}
