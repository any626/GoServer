package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	"github.com/gorilla/securecookie"

	"gopkg.in/yaml.v2"
)

// KataDB is the injected database object
type KataDB struct {
	db *sql.DB
	sc *securecookie.SecureCookie
}

// DbConfig is the database configuration format
type DbConfig struct {
	Username string
	Password string
	DbName   string
}

// Connect to database
func Connect() (KataDB, error) {
	// read credentials from config file
	db := KataDB{}
	d := DbConfig{}
	b, err := ioutil.ReadFile("dbconfig.yaml")
	if err != nil {
		return db, err
	}
	err = yaml.Unmarshal(b, &d)
	if err != nil {
		return db, err
	}
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		d.Username, d.Password, d.DbName)
	// open the db
	db.db, err = sql.Open("postgres", dbinfo)
	return db, err
}

// Disconnect the database
func (db *KataDB) Disconnect() error {
	return db.db.Close()
}
