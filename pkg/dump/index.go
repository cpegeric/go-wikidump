package dump

import (
	"bufio"
	"compress/bzip2"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/BehzadE/go-wikidump/pkg/model"
)

type indexLine struct {
	pageID, byteBegin int64
}

// Get the index files from the dump directory.
func (dump *dump) getIndexFiles() ([]string, error) {
	pattern, err := regexp.Compile(".*index.*")
	result := make([]string, 0)
	if err != nil {
		return nil, err
	}
	err = filepath.Walk(dump.Directory, func(_ string, info os.FileInfo, err error) error {
		if err == nil && pattern.MatchString(info.Name()) {
			result = append(result, dump.Directory+info.Name())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func parseIndexLine(line string) (*indexLine, error) {
	splits := strings.Split(line, ":")
	byteBegin, err := strconv.ParseInt(splits[0], 10, 64)
	if err != nil {
		return nil, err
	}
	pageID, err := strconv.ParseInt(splits[0], 10, 64)
	if err != nil {
		return nil, err
	}
	return &indexLine{pageID: pageID, byteBegin: byteBegin}, nil
}

// bool value shows whether the scanner has reached the end of the file or not.
func scanStream(scanner *bufio.Scanner) ([]*indexLine, bool, error) {
	lines := make([]*indexLine, 0)
	var i int
	for i < 100 && scanner.Scan() {
		i++
		line, err := parseIndexLine(scanner.Text())
		if err != nil {
			return nil, false, err
		}
		lines = append(lines, line)
	}
	return lines, i < 99, nil
}

func parseIndexFile(path string, fileID int64) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	br := bzip2.NewReader(file)
	scanner := bufio.NewScanner(br)
	var prevStreamID, streamID int64
	var done bool
	var streamLines []*indexLine
	for {
		streamLines, done, err = scanStream(scanner)
		if err != nil {
			return err
		}
		streamID, err = model.CreateStream(streamLines[0].byteBegin, fileID)
		if prevStreamID != 0 {
			model.SetStreamByteEnd(prevStreamID, streamLines[0].byteBegin)
		}
		prevStreamID = streamID
		if err != nil {
			return err
		}
		for _, line := range streamLines {
			model.CreatePage(line.pageID, streamID)
		}
		if done {
			break
		}
	}
	model.SetStreamByteEnd(prevStreamID, 0)
	return nil
}

// Read all the index files and store the page byte location and file names to a sqlite database
// for faster querying.
func (dump *dump) ParseIndexes() error {
	indexes, err := dump.getIndexFiles()
	if err != nil {
		return err
	}
	for _, index := range indexes {
		datafile := strings.Replace(index, "txt", "xml", 1)
		datafile = strings.Replace(datafile, "-index", "", 1)
		fileID, err := model.CreateFile(datafile)
		if err != nil {
			return err
		}
		parseIndexFile(index, fileID)
	}
	return nil
}
