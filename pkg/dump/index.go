package dump

import (
	"bufio"
	"compress/bzip2"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/BehzadE/go-wikidump/pkg/model"
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

// bool value shows whether the scanner has reached the end of the file or not.
func readStream(scanner *bufio.Scanner) (int64, []int64, bool, error) {
	pageIDs := make([]int64, 0)
	scanner.Scan()
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

func (d *dump) indexDatafile(datafile *model.Datafile) error {
	fmt.Printf("Parsing index file: %v with fileID: %v\n", datafile.IndexPath, datafile.ID)
	t := time.Now()
	file, err := os.Open(filepath.Join(d.dir, datafile.IndexPath))
	if err != nil {
		return err
	}
	defer file.Close()

	br := bzip2.NewReader(file)
	scanner := bufio.NewScanner(br)
	done := false
	var pageIDs, prevPageIDs []int64
	var byteBegin, prevByteBegin int64
	for !done {
		byteBegin, pageIDs, done, err = readStream(scanner)
		if err != nil {
			return err
		}
		if prevPageIDs != nil {
			err = d.insertStream(datafile.ID, prevByteBegin, byteBegin, prevPageIDs)
			if err != nil {
				return err
			}
		}
		prevByteBegin, prevPageIDs = byteBegin, pageIDs
	}
	if err = d.insertStream(datafile.ID, prevByteBegin, datafile.Size, prevPageIDs); err != nil {
		return err
	}
	d.MarkDatafileIndexed(datafile.ID)
	fmt.Println(time.Since(t))
	return err
}

func (d *dump) initDatafiles() error {
	indexPaths, err := d.getIndexFiles()
	if err != nil {
		return err
	}
	paths := make([]string, len(indexPaths))
	sizes := make([]int64, len(indexPaths))
	for i := range indexPaths {
		path := strings.Replace(indexPaths[i], "txt", "xml", 1)
		path = strings.Replace(path, "-index", "", 1)
		var f *os.File
		var fi os.FileInfo
		f, err = os.Open(filepath.Join(d.dir, path))
		if err != nil {
			return err
		}
		fi, err = f.Stat()
		if err != nil {
			return err
		}
		paths[i] = path
		sizes[i] = fi.Size()
	}
	err = d.insertDatafiles(paths, sizes, indexPaths)
	return err
}

func (d *dump) SaveIndexes() error {
	if err := d.initDatafiles(); err != nil {
		return err
	}
	datafiles, err := d.selectDatafiles()
	if err != nil {
		return err
	}
	for _, datafile := range datafiles {
		if err = d.indexDatafile(datafile); err != nil {
			return err
		}
	}
	return nil
}
