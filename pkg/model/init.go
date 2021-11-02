package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

var migrateUp string = `
create table if not exists pages (
    id integer not null primary key,
    streamid integer not null,
    foreign key(streamid) references streams(id)
);

create table if not exists streams {
    id integer primary key autoincrement,
    beginbyte integer,
    endbyte integer,
    fileid integer not null,
    foreign key(fileid) references files(id)
}

create table if not exists datafiles (
    id integer primary key autoincrement,
    path text,
);
`

func InitDB(dst string) error {
	var err error
	db, err = sql.Open("sqlite3", dst+"db.sqlite3?_foreign_keys=on")
	if err != nil {
		return err
	}
	_, err = db.Exec(migrateUp)
	return err
}
