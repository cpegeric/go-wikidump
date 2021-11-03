package model

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

// Create a row in files table and return the ID.
func CreateFile(db *sql.DB, path string) (int64, error) {
	query := sq.Insert("datafiles").Columns("path").Values(path)
	result, err := query.RunWith(db).Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

func CreateStream(db *sql.DB, byteBegin, fileID int64) (int64, error) {
	query := sq.Insert("streams").Columns("bytebegin", "fileid").Values(byteBegin, fileID)
	result, err := query.RunWith(db).Exec()
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

func SetStreamByteEnd(db *sql.DB, id, byteEnd int64) error {
	query := sq.Update("streams").Set("byteEnd", byteEnd).Where(sq.Eq{"id": id})
	_, err := query.RunWith(db).Exec()
	return err
}

func CreatePage(db *sql.DB, id, streamID int64) error {
	query := sq.Insert("pages").Columns("id", "streamid").Values(id, streamID)
	_, err := query.RunWith(db).Exec()
	return err
}
