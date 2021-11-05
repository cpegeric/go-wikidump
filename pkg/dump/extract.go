package dump

import (
	"bytes"
	"compress/bzip2"
	"io"
	"io/fs"
	"log"
	"os"

	"github.com/BehzadE/go-wikidump/internal/model"
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

func ExtractPages(stream []byte) {
	file, err := os.Open("test.xml")
	if err != nil {
		log.Fatal(err)
	}
	file.Write(stream)
}

// Each stream contains 100 articles. This function parses the xml to get the specified page.
// func getPageFromStream(stream []byte, pageID int64) (*Page, error) {
// 	var pages Pages
// 	buff := bytes.NewBufferString("<pages>\n")
// 	buff.Write(stream)
// 	buff.WriteString("</pages>")
// 	err := xml.Unmarshal(buff.Bytes(), &pages)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	for _, page := range pages.Pages {
// 		if page.ID == pageID {
// 			return &page, nil
// 		}
// 	}
// 	return nil, errors.New("not found")
// }
//
// // Finds and returns the page associated with the specified pageID if it exists.
// func (dump *Dump) GetPage(pageID int64) (*Page, error) {
// 	indexFile, err := dump.findIndex(pageID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	byteLocations, err := getPageByteLocation(indexFile, pageID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	stream, err := getStream(indexFile, byteLocations)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	page, err := getPageFromStream(stream, pageID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return page, nil
// }
