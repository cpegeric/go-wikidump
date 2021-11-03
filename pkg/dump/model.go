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

func (d *dump) createFile(path string) (int64, error) {
	return model.CreateFile(d.db, path)
}

func (d *dump) createStream(bytebegin, fileID int64) (int64, error) {
	return model.CreateStream(d.db, bytebegin, fileID)
}

func (d *dump) setStreamByteEnd(id, byteEnd int64) error {
	return model.SetStreamByteEnd(d.db, id, byteEnd)
}

func (d *dump) createPage(id, streamID int64) error {
	return model.CreatePage(d.db, id, streamID)
}

func (d *dump) getStreamID(pageID int64) (int64, error) {
	return model.GetStreamID(d.db, pageID)
}

func (d *dump) getStream(streamID int64) (*model.Stream, error) {
	return model.GetStream(d.db, streamID)
}
