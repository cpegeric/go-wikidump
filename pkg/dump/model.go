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

func (d *dump) insertStream(archiveID, byteBegin, byteEnd int64, pageIDs []int64) error {
	return model.InsertStream(d.db, archiveID, byteBegin, byteEnd, pageIDs)
}

func (d *dump) selectPage(pageID int64) (int64, error) {
	return model.SelectPage(d.db, pageID)
}

func (d *dump) selectStream(streamID int64) (*model.Stream, error) {
	return model.SelectStream(d.db, streamID)
}

func (d *dump) insertArchives(archives []*model.Archive) error {
	return model.InsertArchives(d.db, archives)
}

func (d *dump) selectArchives() ([]*model.Archive, error) {
	return model.SelectArchives(d.db)
}

func (d *dump) markArchiveProcessed(id int64) error {
	return model.MarkArchiveProcessed(d.db, id)
}
