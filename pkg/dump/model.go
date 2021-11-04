package dump

import (
	"database/sql"
	"path"

	"github.com/BehzadE/go-wikidump/pkg/model"
)

type dump struct {
	db  *sql.DB
	dir string
}

func New(dir string) (*dump, error) {
	dir = path.Clean(dir)
	db, err := model.InitDB(dir)
	if err != nil {
		return nil, err
	}
	d := dump{dir: dir, db: db}
	return &d, nil
}

func (d *dump) insertStream(fileID, byteBegin, byteEnd int64, pageIDs []int64) error {
	return model.InsertStream(d.db, fileID, byteBegin, byteEnd, pageIDs)
}

func (d *dump) getStreamID(pageID int64) (int64, error) {
	return model.SelectPageStreamID(d.db, pageID)
}

func (d *dump) getStream(streamID int64) (*model.Stream, error) {
	return model.SelectStreamByID(d.db, streamID)
}

func (d *dump) insertDatafiles(datafiles []string, sizes []int64, indexfiles []string) error {
	return model.InsertDatafiles(d.db, datafiles, sizes, indexfiles)
}

func (d *dump) selectDatafiles() ([]*model.Datafile, error) {
	return model.SelectDatafiles(d.db)
}

func (d *dump) MarkDatafileIndexed(id int64) error {
	return model.MarkDatafileIndexed(d.db, id)
}
