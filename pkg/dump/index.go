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

func (d *dump) processArchiveIndex(archive *model.Archive) error {
	fmt.Printf("Index file: %v with : %v\n", archive.IndexPath, archive.ID)
	t := time.Now()
	file, err := os.Open(filepath.Join(d.dir, archive.IndexPath))
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
			err = d.insertStream(archive.ID, prevByteBegin, byteBegin, prevPageIDs)
			if err != nil {
				return err
			}
		}
		prevByteBegin, prevPageIDs = byteBegin, pageIDs
	}
	if err = d.insertStream(archive.ID, prevByteBegin, archive.FileSize, prevPageIDs); err != nil {
		return err
	}
	d.markArchiveProcessed(archive.ID)
	fmt.Println(time.Since(t))
	return err
}

func (d *dump) saveArchives() error {
	indexPaths, err := d.getIndexFiles()
	if err != nil {
		return err
	}
	var archives []*model.Archive
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
		archive := model.Archive{FilePath: path, FileSize: fi.Size(), IndexPath: indexPaths[i]}
		archives = append(archives, &archive)
	}
	err = d.insertArchives(archives)
	return err
}

func (d *dump) SaveIndexes() error {
	if err := d.saveArchives(); err != nil {
		return err
	}
	archives, err := d.selectArchives()
	if err != nil {
		return err
	}
	for _, archive := range archives {
		if err = d.processArchiveIndex(archive); err != nil {
			return err
		}
	}
	return nil
}
