package database

import (
	"fmt"

	"github.com/gorilla/securecookie"
)

// GenerateSecureCookie gets new crypto
func (db *KataDB) GenerateSecureCookie() {
	row := db.db.QueryRow("SELECT hash,block FROM securecookie LIMIT 1")
	var blockKey, hashKey []byte
	err := row.Scan(&hashKey, &blockKey)
	if err != nil {
		fmt.Println("SecureCookie not found, generating new SecureCookie.")
		fmt.Println("All users will need to refresh their logins.")
		hashKey = securecookie.GenerateRandomKey(64)
		blockKey = securecookie.GenerateRandomKey(32)
		db.db.QueryRow("INSERT INTO securecookie(hash,block) VALUES($1,$2)", hashKey, blockKey)
	}
	db.sc = securecookie.New(hashKey, blockKey)
}

// SecureEncode encodes with our secure cookie
func (db *KataDB) SecureEncode(name string, value interface{}) (string, error) {
	return db.sc.Encode(name, value)
}

// SecureDecode decodes with our secure cookie
func (db *KataDB) SecureDecode(name, value string, dest interface{}) error {
	return db.sc.Decode(name, value, dest)
}
