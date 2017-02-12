package database

import (
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

// SecureEncode encodes with our secure cookie
func SecureEncode(name string, value interface{}) (string, error) {
	return sc.Encode(name, value)
}

// SecureDecode decodes with our secure cookie
func SecureDecode(name, value string, dest interface{}) error {
	return sc.Decode(name, value, dest)
}
