package database

import (
	"net/http"

	"fmt"

	"github.com/gorilla/securecookie"
)

var sc *securecookie.SecureCookie

// GenerateSecureCookie gets new crypto
func GenerateSecureCookie() {
	row := db.QueryRow("SELECT hash,block FROM securecookie LIMIT 1")
	var blockKey, hashKey []byte
	err := row.Scan(&hashKey, &blockKey)
	if err != nil {
		fmt.Println("err not nil")
		hashKey = securecookie.GenerateRandomKey(64)
		blockKey = securecookie.GenerateRandomKey(32)
		db.QueryRow("INSERT INTO securecookie(hash,block) VALUES($1,$2)", hashKey, blockKey)
	}
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
