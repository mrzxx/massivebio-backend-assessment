package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	DB_HOST     = "localhost"
	DB_USER     = "postgres"
	DB_PASSWORD = "admin"
	DB_NAME     = "massivebio"
)

var db *sql.DB

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		fmt.Println("Could not connect to the Postgres Database")
		panic(err)
	}
	return db
}

func GetDatabase() *sql.DB {
	return db
}
