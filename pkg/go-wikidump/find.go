package gowikidump

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"
)

// Find the index file for a given pageID from the offsets calculted before.
func (dump *Dump) findIndex(pageID int64) (string, error) {
	file, err := os.Open(dump.Parameters.DumpDirectory + "offsets.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), "###")
		offset, filename := row[0], row[1]
		offsetNum, err := strconv.Atoi(offset)
		if err != nil {
			return "", err
		}

		if int64(offsetNum) >= pageID {
			return filename, nil
		}
	}
	return "", errors.New("not found")
}

// Returns the position of the stream containing the article in the bzip2 file. If byteEnd is zero
// it indicates that the endByte is the end of the file.
func getPageByteLocation(filename string, pageID int64) ([]int64, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	br := bzip2.NewReader(file)
	scanner := bufio.NewScanner(br)

	var line string
	var found bool
	var byteBegin int64
	var byteEnd int64
	var lineCounter int
	for scanner.Scan() {
		line = scanner.Text()
		if !found && strings.Contains(line, ":"+fmt.Sprint(pageID)+":") {
			split := strings.Split(line, ":")
			byteBegin, err = strconv.ParseInt(split[0], 10, 64)
			if err != nil {
				return nil, err
			}
			found = true
		}
		if found {
			lineCounter++
		}
		if lineCounter > 100 {
			split := strings.Split(line, ":")
			byteEnd, err = strconv.ParseInt(split[0], 10, 64)
			if err != nil {
				return nil, err
			}
			break
		}
	}
	return []int64{byteBegin, byteEnd}, nil
}

// Get the stream containing the page from bz2 file.
func getStream(indexFilename string, byteLocations []int64) ([]byte, error) {
	dataFilename := strings.Replace(indexFilename, "txt", "xml", 1)
	dataFilename = strings.Replace(dataFilename, "-index", "", 1)
	file, err := os.Open(dataFilename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var isEnd bool
	if byteLocations[1] == 0 {
		var fi fs.FileInfo
		fi, err = file.Stat()
		if err != nil {
			return nil, err
		}
		byteLocations[1] = fi.Size()
		isEnd = true
	}

	sr := io.NewSectionReader(file, byteLocations[0], byteLocations[1]-byteLocations[0])
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

// Each stream contains 100 articles. This function parses the xml to get the specified page.
func getPageFromStream(stream []byte, pageID int64) (*Page, error) {
	var pages Pages
	buff := bytes.NewBufferString("<pages>\n")
	buff.Write(stream)
	buff.WriteString("</pages>")
	err := xml.Unmarshal(buff.Bytes(), &pages)
	if err != nil {
		return nil, err
	}

	for _, page := range pages.Pages {
		if page.ID == pageID {
			return &page, nil
		}
	}
	return nil, errors.New("not found")
}

// Finds and returns the page associated with the specified pageID if it exists.
func (dump *Dump) GetPage(pageID int64) (*Page, error) {
	indexFile, err := dump.findIndex(pageID)
	if err != nil {
		return nil, err
	}
	byteLocations, err := getPageByteLocation(indexFile, pageID)
	if err != nil {
		return nil, err
	}

	stream, err := getStream(indexFile, byteLocations)
	if err != nil {
		return nil, err
	}

	page, err := getPageFromStream(stream, pageID)
	if err != nil {
		return nil, err
	}
	return page, nil
}
