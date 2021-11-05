package model

import (
	"database/sql"
	"fmt"
)

// Bz2 files are split into streams. Each stream contains 100 articles. A stream
// is the smallest block we can extract.
type Stream struct {
	Path      string
	ByteBegin int64
	ByteEnd   int64
}

func InsertStream(db *sql.DB, archiveID, byteBegin, byteEnd int64, pageIDs []int64) error {
	_, err := streamExists(db, byteBegin, archiveID)
	if err == nil {
		return nil
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	query := `
        insert or ignore
        into Stream (ArchiveID,ByteBegin,ByteEnd) 
        values(?,?,?)
    `
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(archiveID, byteBegin, byteEnd)
	if err != nil {
		return err
	}
	streamID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	query = `
        insert or ignore
        into Page (ID,StreamID) 
        values 
    `
	for _, pageID := range pageIDs {
		query += fmt.Sprintf("(%v,%v),", pageID, streamID)
	}
	query = query[:len(query)-1] + ";"
	stmt, err = tx.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func SelectStream(db *sql.DB, streamID int64) (*Stream, error) {
	query := `
        select s.ByteBegin, s.ByteEnd, a.Path
        from Stream s 
        where s.ID = ?
        inner join Archive a
        on s.ArchiveID = a.ID
    `
	var stream Stream
	err := db.QueryRow(query, streamID).Scan(&stream.ByteBegin, &stream.ByteEnd, &stream.Path)
	if err != nil {
		return nil, err
	}
	return &stream, nil
}

func streamExists(db *sql.DB, byteBegin, archiveID int64) (int64, error) {
	query := `
        select ID
        from Stream
        where ByteBegin = ? and ArchiveID = ?

    `
	var streamID int64
	err := db.QueryRow(query, byteBegin, archiveID).Scan(&streamID)
	return streamID, err
}
