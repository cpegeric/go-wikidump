package gowikidump

import "encoding/xml"

type Index struct {
	Filename string `json:"filename"`
	EndID    string `json:"endid"`
}

type Page struct {
	XMLName   xml.Name `xml:"page"`
	Title     string   `xml:"title"`
	Namespace string   `xml:"ns"`
	ID        int64    `xml:"id"`
	Revision  Revision `xml:"revision"`
}

type Revision struct {
	XMLName xml.Name `xml:"revision"`
	ID      int64    `xml:"id"`
	Format  string   `xml:"format"`
	Text    string   `xml:"text"`
}

type Pages struct {
	XMLName xml.Name `xml:"pages"`
	Pages   []Page   `xml:"page"`
}

type Parameters struct {
	BaseURL       string `default:"https://dumps.wikimedia.org"`
	DumpVer       string `default:"/enwiki/20210720/"`
	DumpDirectory string `default:"./wikipedia-dump/"`
}

type Dump struct {
	Parameters Parameters
	Links      []string
}
