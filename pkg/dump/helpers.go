package dump

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Get the index files from the dump directory.
func (dump *dump) getIndexFiles() ([]string, error) {
	pattern, err := regexp.Compile(".*index.*")
	result := make([]string, 0)
	if err != nil {
		return nil, err
	}
	err = filepath.Walk(dump.dir, func(_ string, info os.FileInfo, err error) error {
		if err == nil && pattern.MatchString(info.Name()) {
			result = append(result, info.Name())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Parse a single line in a index file.
func parseIndexLine(line string) (int64, int64, error) {
	splits := strings.Split(line, ":")
	byteBegin, err := strconv.ParseInt(splits[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	pageID, err := strconv.ParseInt(splits[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return byteBegin, pageID, nil
}

type scannerExhaustedError struct{}

func (err scannerExhaustedError) Error() string {
	return "no more lines in scanner"
}

// bool value shows whether the scanner has reached the end of the file or not.
func readStream(scanner *bufio.Scanner) (int64, []int64, bool, error) {
	pageIDs := make([]int64, 0)
	if !scanner.Scan() {
		return 0, nil, true, scannerExhaustedError{}
	}
	byteBegin, pageID, err := parseIndexLine(scanner.Text())
	if err != nil {
		return 0, nil, false, err
	}
	pageIDs = append(pageIDs, pageID)
	i := 1
	for i < 100 && scanner.Scan() {
		i++
		_, pageID, err = parseIndexLine(scanner.Text())
		if err != nil {
			return 0, nil, false, err
		}
		pageIDs = append(pageIDs, pageID)
	}
	return byteBegin, pageIDs, i < 99, nil
}
