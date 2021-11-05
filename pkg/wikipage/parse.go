package wikipage

import (
	"bytes"
	"encoding/xml"
	"time"
)

type Stream struct {
	XMLName xml.Name `xml:"stream"`
	Pages   []*Page  `xml:"page"`
}

type Page struct {
	Title    string   `xml:"title"`
	Redirect Redirect `xml:"redirect"`
	Revision Revision `xml:"revision"`
	NS       int64    `xml:"ns"`
	ID       int64    `xml:"id"`
}

type Redirect struct {
	Title string `xml:"title,attr"`
}

type Revision struct {
	Timestamp   time.Time   `xml:"timestamp"`
	Format      string      `xml:"format"`
	Text        string      `xml:"text"`
	Comment     string      `xml:"comment"`
	Model       string      `xml:"model"`
	SHA1        string      `xml:"sha1"`
	Contributer Contributer `xml:"contributer"`
	ID          int64       `xml:"id"`
	ParentID    int64       `xml:"parentid"`
}

type Contributer struct {
	Username string `xml:"username"`
	ID       int64  `xml:"id"`
}

func ParseStream(stream []byte) ([]*Page, error) {
	var s Stream
	buff := bytes.NewBufferString("<stream>\n")
	buff.Write(stream)
	buff.WriteString("</stream>")
	err := xml.Unmarshal(buff.Bytes(), &s)
	if err != nil {
		return nil, err
	}
	return s.Pages, nil
}

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
