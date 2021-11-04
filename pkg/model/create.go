package model

import (
	"database/sql"
	"fmt"
)

func InsertDatafiles(db *sql.DB, datafiles []string, sizes []int64, indexfiles []string) error {
	query := "insert or ignore into datafiles (path,size,indexpath) values "
	for i := range datafiles {
		query += fmt.Sprintf("('%v',%v,'%v'),", datafiles[i], sizes[i], indexfiles[i])
	}
	query = query[:len(query)-1] + ";"
	_, err := db.Exec(query)
	return err
}

func MarkDatafileIndexed(db *sql.DB, id int64) error {
	query := "update datafiles set indexed = true where id = ?"
	_, err := db.Exec(query, id)
	return err
}

func InsertStream(db *sql.DB, fileID, byteBegin, byteEnd int64, pageIDs []int64) error {
	_, err := selectStream(db, byteBegin, fileID)
	if err == nil {
		return nil
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	query := "insert or ignore into streams (fileid,bytebegin,byteend) values(?,?,?)"
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	result, err := stmt.Exec(fileID, byteBegin, byteEnd)
	if err != nil {
		return err
	}
	streamID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	query = "insert or ignore into pages (id,streamid) values "
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
