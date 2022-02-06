package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
)

var dbUser string
var dbPassword string
var dbHost string
var dbPort int
var dbName string
var dbTZ string

func init() {
	dbUser = os.Getenv("DB_USER")
	if dbUser == "" {
		panic(errors.New("DB_USER is unset"))
	}
	dbPassword = os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		panic(errors.New("DB_PASSWORD is unset"))
	}
	dbHost = os.Getenv("DB_HOST")
	if dbHost == "" {
		panic(errors.New("DB_HOST is unset"))
	}
	dbPortStr := os.Getenv("DB_PORT")
	if dbPortStr == "" {
		panic(errors.New("DB_PORT is unset"))
	} else {
		var e error
		dbPort, e = strconv.Atoi(dbPortStr)
		if e != nil {
			panic(e)
		}
	}
	dbName = os.Getenv("DB_NAME")
	if os.Getenv("DB_NAME") == "" {
		panic(errors.New("DB_NAME is unset"))
	}
	dbTZ = os.Getenv("TZ")
	if os.Getenv("TZ") == "" {
		panic(errors.New("TZ is unset"))
	}
}

func NewDB() (*sql.DB, error) {
	db, e := sql.Open("mysql", fmt.Sprintf("%v:%v@tcp(%v:%d)/%v?parseTime=true&loc=%v", dbUser, dbPassword, dbHost, dbPort, dbName, url.QueryEscape(dbTZ)))
	if e != nil {
		return nil, e
	}

	db.SetMaxIdleConns(1) // retain a connection to maintain the session between queries in a request
	db.SetMaxOpenConns(1)
	return db, nil
}
