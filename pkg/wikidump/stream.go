package wikidump

import (
	"bytes"
	"compress/bzip2"
	"io"
	"os"
	"path/filepath"

	"github.com/BehzadE/go-wikidump/internal/model"
)

// Get the stream containing the page from bz2 file.
func extractStream(path string, byteBegin, byteEnd int64) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	sr := io.NewSectionReader(file, byteBegin, byteEnd-byteBegin)
	reader := bzip2.NewReader(sr)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	// In the case where the stream is at the end of the file strip the extra closing </mediawiki> tag.
	if byteEnd == fi.Size() {
		index := bytes.LastIndex(data, []byte("\n"))
		data = data[:index]
	}
	return data, nil
}

type streamReader struct {
	path    string
	streams []*model.Stream
	pointer int
}

func (sr *streamReader) Next() bool {
	sr.pointer++
	return sr.pointer < len(sr.streams)
}

func (sr *streamReader) Read() ([]byte, error) {
	stream := sr.streams[sr.pointer]
	return extractStream(sr.path, stream.ByteBegin, stream.ByteEnd)
}

func (d *dump) NewStreamReader(archivePath string) (*streamReader, error) {
	streams, err := d.selectArchiveStreams(archivePath)
	if err != nil {
		return nil, err
	}
	return &streamReader{
		streams: streams,
		pointer: -1,
		path:    filepath.Join(d.dir, archivePath),
	}, nil
}
