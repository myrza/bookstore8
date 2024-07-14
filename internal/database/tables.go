package server

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func TableCreate(db *sql.DB) error {
	// create table if not exists
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS authors (id SERIAL PRIMARY KEY, name TEXT, surname TEXT, biography TEXT, birthday DATE)")
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS books (id SERIAL PRIMARY KEY, title TEXT, authorid INTEGER, isbn TEXT, year INTEGER)")
	if err != nil {
		return err
	}
	return nil
}
