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

func GetStreamID(db *sql.DB, pageID int64) (int64, error) {
	query := sq.Select("p.streamid").
		From("pages p").
		Where(sq.Eq{"p.id": pageID})
	var streamID int64
	err := query.RunWith(db).QueryRow().Scan(&streamID)
	if err != nil {
		return 0, err
	}
	return streamID, nil
}

func GetStream(db *sql.DB, streamID int64) (*Stream, error) {
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
