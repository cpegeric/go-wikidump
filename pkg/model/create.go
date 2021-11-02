package model

import sq "github.com/Masterminds/squirrel"

// Create a row in files table and return the ID.
func CreateFile(path string) (int64, error) {
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

func CreateStream(beginByte, fileID int64) (int64, error) {
	query := sq.Insert("streams").Columns("beginbyte", "fileid").Values(beginByte, fileID)
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

func SetStreamByteEnd(id, byteEnd int64) error {
	query := sq.Update("streams").Set("byteEnd", byteEnd).Where(sq.Eq{"id": id})
	_, err := query.RunWith(db).Exec()
	return err
}

func CreatePage(id, streamID int64) error {
	query := sq.Insert("pages").Columns("id", "streamid").Values(id, streamID)
	_, err := query.RunWith(db).Exec()
	return err
}
