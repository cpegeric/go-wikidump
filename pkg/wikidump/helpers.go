package wikidump

import (
	"fmt"
	"bufio"
	"strconv"
	"strings"
)

type indexline struct {
	pageName  string
	byteBegin int64
	pageID    int64
}

type stream struct {
	pageIDs   []int64
	pageNames []string
	byteBegin int64
	last      bool
}

// Parse a single line in a index file.
func splitIndexline(line string) (*indexline, error) {
	splits := strings.SplitN(line, ":", 3)
	var result indexline
	var err error
	result.byteBegin, err = strconv.ParseInt(splits[0], 10, 64)
	if err != nil {
		return nil, err
	}
	result.pageID, err = strconv.ParseInt(splits[1], 10, 64)
	if err != nil {
		return nil, err
	}
	result.pageName = splits[2]
	return &result, nil
}

type scannerExhaustedError struct{}

func (err scannerExhaustedError) Error() string {
	return "no more lines in scanner"
}

// Read a hundred lines from the index file or until the end of file.
// bool value shows whether the scanner has reached the end of the file or not.
func readStream(scanner *bufio.Scanner) (*stream, error) {
	var s stream
	fmt.Printf("read stream start")
	if !scanner.Scan() {
		return nil, scannerExhaustedError{}
	}
	il, err := splitIndexline(scanner.Text())
	if err != nil {
		return nil, err
	}
	s.byteBegin = il.byteBegin
	s.pageNames = append(s.pageNames, il.pageName)
	s.pageIDs = append(s.pageIDs, il.pageID)
	i := 1
	for i < 100 && scanner.Scan() {
		i++
		il, err = splitIndexline(scanner.Text())
		if err != nil {
			return nil, err
		}
		s.pageNames = append(s.pageNames, il.pageName)
		s.pageIDs = append(s.pageIDs, il.pageID)
	}
	fmt.Printf("readStream end... %d", i)
	s.last = i < 99
	return &s, nil
}
