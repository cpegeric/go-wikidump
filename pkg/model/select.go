package model

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

// Bz2 files are split into streams. Each stream contains 100 articles. A stream
// is the smallest block we can extract.
type Stream struct {
	Path      string
	ByteBegin int64
	ByteEnd   int64
}

type Datafile struct {
	Path      string
	IndexPath string
	ID        int64
	Size      int64
}

func SelectPageStreamID(db *sql.DB, pageID int64) (int64, error) {
	query := "select streamid from pages where id = ?"
	var streamID int64
	err := db.QueryRow(query, pageID).Scan(&streamID)
	if err != nil {
		return 0, err
	}
	return streamID, nil
}

func SelectStreamByID(db *sql.DB, streamID int64) (*Stream, error) {
	query := sq.Select("s.bytebegin", "s.byteend", "f.path").
		From("streams s").
		Where(sq.Eq{"s.id": streamID}).
		InnerJoin("datafiles f on s.fileid=f.id")
	var stream Stream
	err := query.RunWith(db).QueryRow().Scan(&stream.ByteBegin, &stream.ByteEnd, &stream.Path)
	if err != nil {
		return nil, err
	}
	return &stream, nil
}

func SelectDatafiles(db *sql.DB) ([]*Datafile, error) {
	query := "select id,path,size,indexpath from datafiles where indexed = false"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	var results []*Datafile
	for rows.Next() {
		var df Datafile
		if err := rows.Scan(&df.ID, &df.Path, &df.Size, &df.IndexPath); err != nil {
			return nil, err
		}
		results = append(results, &df)
	}
	return results, nil
}

func selectStream(db *sql.DB, byteBegin, fileID int64) (int64, error) {
	query := sq.Select("id").From("streams").Where(sq.Eq{"bytebegin": byteBegin}, sq.Eq{"fileid": fileID})
	var streamID int64
	err := query.RunWith(db).QueryRow().Scan(&streamID)
	return streamID, err
}
