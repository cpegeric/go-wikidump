package dump

import (
	"bytes"
	"compress/bzip2"
	"io"
	"io/fs"
	"os"

	"github.com/BehzadE/go-wikidump/pkg/model"
)

// Get the stream containing the page from bz2 file.
func ExtractStream(streamInfo *model.Stream) ([]byte, error) {
	file, err := os.Open(streamInfo.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var isEnd bool
	if streamInfo.ByteEnd == 0 {
		var fi fs.FileInfo
		fi, err = file.Stat()
		if err != nil {
			return nil, err
		}
		streamInfo.ByteEnd = fi.Size()
		isEnd = true
	}

	sr := io.NewSectionReader(file, streamInfo.ByteBegin, streamInfo.ByteEnd-streamInfo.ByteBegin)
	reader := bzip2.NewReader(sr)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	// In the case where the stream is at the end of the file strip the extra closing </mediawiki> tag.
	if isEnd {
		index := bytes.LastIndex(data, []byte("\n"))
		data = data[:index]
	}
	return data, nil
}

func ExtractPages(stream []byte) []*Page
