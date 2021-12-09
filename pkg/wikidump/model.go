package wikidump

import (
	"database/sql"
	"path"

	"github.com/BehzadE/go-wikidump/internal/model"
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

func (d *dump) insertStream(archiveID, byteBegin, byteEnd int64, pageIDs []int64, pageNames []string) error {
	return model.InsertStream(d.db, archiveID, byteBegin, byteEnd, pageIDs, pageNames)
}

func (d *dump) selectPages(pageIDs []int64) ([]int64, error) {
	return model.SelectPages(d.db, pageIDs)
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

func (d *dump) selectArchiveStreams(archivePath string) ([]*model.Stream, error) {
	return model.SelectArchiveStreams(d.db, archivePath)
}

func (d *dump) getPageStreams(pageIDs []int64) (map[*model.Stream][]int64, error) {
	streamIDs, err := model.SelectPages(d.db, pageIDs)
	if err != nil {
		return nil, err
	}
	streamPage := make(map[int64][]int64)
	for i, streamID := range streamIDs {
		if streamID != 0 {
			streamPage[streamID] = append(streamPage[streamID], pageIDs[i])
		}
	}
	unique := []int64{}
	for k := range streamPage {
		unique = append(unique, k)
	}
	streams, err := model.SelectStreams(d.db, unique)
	if err != nil {
		return nil, err
	}
	result := make(map[*model.Stream][]int64)
	for _, stream := range streams {
		result[stream] = streamPage[stream.ID]
	}
	return result, nil
}
