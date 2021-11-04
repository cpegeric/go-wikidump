package model

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dst string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dst+"/db.sqlite3?_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	query, err := os.ReadFile("createDB.sql")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(string(query))
	if err != nil {
		return nil, err
	}
	return db, nil
}
