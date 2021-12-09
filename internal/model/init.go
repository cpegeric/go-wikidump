package model

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var migration = `
create table if not exists Archive (
    ID integer primary key autoincrement,
    FilePath text unique,
    FileSize integer,
    IndexPath text unique,
    Processed boolean default false
);

create table if not exists Page (
    ID integer primary key,
    Name text,
    StreamID integer not null,
    foreign key(StreamID) references Stream(id)
);

create table if not exists Stream (
    ID integer primary key autoincrement,
    ByteBegin integer,
    ByteEnd integer,
    ArchiveID integer not null,
    foreign key(ArchiveID) references Archive(ID),
    unique(ByteBegin,ArchiveID)
);
`

func InitDB(dst string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dst+"/db.sqlite3?_foreign_keys=on&cache=shared")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	_, err = db.Exec(string(migration))
	if err != nil {
		return nil, err
	}
	return db, nil
}
