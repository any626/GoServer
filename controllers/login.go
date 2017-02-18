package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/brwhale/GoServer/database"
	"golang.org/x/crypto/bcrypt"
)

type validation struct {
	CurrentUser string
	Errors      bool
}

// Login is the login page for the site
func (c *KataController) Login(w http.ResponseWriter, r *http.Request) {
	Username := r.FormValue("Username")
	Password := r.FormValue("Password")
	val := validation{Errors: false, CurrentUser: c.GetSecureUsername(r)}
	if Username != "" && Password != "" {
		if c.DB.ValidatePassword(Username, Password) {
			// redirect
			c.FinalizeLogin(Username, w, r)
		}
		val.Errors = true
	}
	c.templates.ExecuteTemplate(w, "Login", val)
}

// Logout is the logout page for the site
func (c *KataController) Logout(w http.ResponseWriter, r *http.Request) {
	c.FinalizeLogin("", w, r)
}

// Register is the register page for the site
func (c *KataController) Register(w http.ResponseWriter, r *http.Request) {
	Username := r.FormValue("Username")
	Password := r.FormValue("Password")
	val := validation{Errors: false, CurrentUser: c.GetSecureUsername(r)}
	if Username != "" && Password != "" {
		/*
			Don't let people register names with these characters:
			/ ? #
			They don't work in URLs for "/user/{username}"
		*/
		if !strings.ContainsAny(Username, "/?#") && len(Password) > 9 {
			hash, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, err.Error(), 500)
			}
			err = c.DB.InsertUser(database.User{Name: Username, Hash: string(hash)})
			if err == nil {
				// redirect
				c.FinalizeLogin(Username, w, r)
			}
		}
		val.Errors = true
	}
	c.templates.ExecuteTemplate(w, "Register", val)
}

// FinalizeLogin generates cookie and save session info
func (c *KataController) FinalizeLogin(username string, w http.ResponseWriter, r *http.Request) {
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
		if encoded, err := c.DB.SecureEncode("login-name", value); err == nil {
			cookie := &http.Cookie{
				Name:  "login-name",
				Value: encoded,
				Path:  "/",
			}
			http.SetCookie(w, cookie)
		}
	}
	http.Redirect(w, r, "/boards", 302)
}

// GetSecureUsername from secure cookie
func (c *KataController) GetSecureUsername(r *http.Request) string {
	Username := ""
	if cookie, err := r.Cookie("login-name"); err == nil {
		value := make(map[string]string)
		if err = c.DB.SecureDecode("login-name", cookie.Value, &value); err == nil {
			Username = value["username"]
		}
	}
	return Username
}
