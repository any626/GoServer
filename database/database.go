package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/brwhale/GoServer/util"

	"gopkg.in/yaml.v2"
)

// database object
var db *sql.DB

// DbConfig is the database configuration format
type DbConfig struct {
	Username string
	Password string
	DbName   string
}

// Connect to database
func Connect() {
	// read credentials from config file
	d := DbConfig{}
	b, err := ioutil.ReadFile("dbconfig.yaml")
	util.Check(err)
	err = yaml.Unmarshal(b, &d)
	util.Check(err)
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		d.Username, d.Password, d.DbName)
	// open the db
	db, err = sql.Open("postgres", dbinfo)
	util.Check(err)
}

// Disconnect the database
func Disconnect() {
	db.Close()
}

func friendlyString(duration time.Duration) string {
	if duration.Hours() >= 48 {
		return fmt.Sprintf("%.0f days ago", duration.Hours()/24)
	}
	if duration.Hours() >= 24 {
		return "1 day ago"
	}
	if duration.Hours() >= 2 {
		return fmt.Sprintf("%.0f hours ago", duration.Hours())
	}
	if duration.Hours() >= 1 {
		return "1 hour ago"
	}
	if duration.Minutes() >= 2 {
		return fmt.Sprintf("%.0f minutes ago", duration.Minutes())
	}
	if duration.Minutes() >= 1 {
		return "1 minute ago"
	}
	return "a couple seconds ago"
}
