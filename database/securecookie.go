package database

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

var sc *securecookie.SecureCookie

// GenerateSecureCookie gets new crypto
func GenerateSecureCookie() {
	var hashKey = securecookie.GenerateRandomKey(64)
	var blockKey = securecookie.GenerateRandomKey(32)
	sc = securecookie.New(hashKey, blockKey)
}

// GetSecureUsername from secure cookie
func GetSecureUsername(r *http.Request) string {
	Username := ""
	if cookie, err := r.Cookie("login-name"); err == nil {
		value := make(map[string]string)
		if err = sc.Decode("login-name", cookie.Value, &value); err == nil {
			Username = value["username"]
		}
	}
	return Username
}

// SecureEncode encodes with our secure cookie
func SecureEncode(name string, value interface{}) (string, error) {
	return sc.Encode(name, value)
}
