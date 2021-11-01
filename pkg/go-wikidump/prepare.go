package gowikidump

import (
	"bufio"
	"compress/bzip2"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Get the index files from the dump directory.
func (dump *Dump) getIndexFiles() ([]string, error) {
	pattern, err := regexp.Compile(".*index.*")
	result := make([]string, 0)
	if err != nil {
		return nil, err
	}
	err = filepath.Walk(dump.Parameters.DumpDirectory, func(_ string, info os.FileInfo, err error) error {
		if err == nil && pattern.MatchString(info.Name()) {
			result = append(result, dump.Parameters.DumpDirectory+info.Name())
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Finds the last ID that is contained within each index file and return the slice.
func (dump *Dump) getIndexRanges() ([]Index, error) {
	indexFiles, err := dump.getIndexFiles()
	if err != nil {
		return nil, err
	}
	indexSlice := make([]Index, 0)
	for _, indexFile := range indexFiles {
		tail, err := getIndexLastLine(indexFile)
		if err != nil {
			return nil, err
		}
		toAdd := Index{
			EndID:    strings.Split(string(tail), ":")[1],
			Filename: indexFile,
		}
		indexSlice = append(indexSlice, toAdd)
	}
	return indexSlice, nil
}

// Helper function to mimic unix's tail -n 1 command but for a bz2 compressed file.
func getIndexLastLine(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	br := bzip2.NewReader(file)
	scanner := bufio.NewScanner(br)

	var line []byte
	for scanner.Scan() {
		line = scanner.Bytes()
	}
	return line, nil
}

// Saves the ID ranges for index files in the offsets.txt file in the dump directory to
// avoid recalculating it in the future.
func (dump *Dump) SaveIndexRanges() error {
	indexRanges, err := dump.getIndexRanges()
	if err != nil {
		return err
	}
	file, err := os.Create(dump.Parameters.DumpDirectory + "offsets.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	sort.Slice(indexRanges, func(i int, j int) bool {
		int1, err := strconv.Atoi(indexRanges[i].EndID)
		if err != nil {
			panic(err)
		}
		int2, err := strconv.Atoi(indexRanges[j].EndID)
		if err != nil {
			panic(err)
		}
		return int1 < int2
	})

	for _, row := range indexRanges {
		fmt.Fprintln(file, row.EndID+"###"+row.Filename)
	}
	return nil
}
