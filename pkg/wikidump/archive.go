package wikidump

import (
	"bufio"
	"compress/bzip2"
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/BehzadE/go-wikidump/internal/model"
)

// Get the index files from the dump directory.
func (d *dump) walkDumpDir() ([]string, error) {
	pattern, err := regexp.Compile(".*index.*")
	result := make([]string, 0)
	if err != nil {
		return nil, err
	}
	err = filepath.Walk(d.dir, func(_ string, info os.FileInfo, err error) error {
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

// Read a single index file and save lines to the database.
func (d *dump) processArchiveIndex(archive *model.Archive) error {
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
			if errors.Is(err, scannerExhaustedError{}) {
				break
			}
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
	err = d.markArchiveProcessed(archive.ID)
	return err
}

// Get archives from database and process those that have not already been processed.
func (d *dump) processArchives() error {
	indexPaths, err := d.walkDumpDir()
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

// Populates the database with the index information available in dump directory.
func (d *dump) PopulateDB() error {
	if err := d.processArchives(); err != nil {
		return err
	}
	archives, err := d.selectArchives()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, archive := range archives {
		wg.Add(1)
		go func(wg *sync.WaitGroup, archive *model.Archive) {
			defer wg.Done()
			if err = d.processArchiveIndex(archive); err != nil {
				log.Printf("error processing archive index with ID %v: %v", archive.ID, err)
			}
		}(&wg, archive)
	}
	wg.Wait()
	return nil
}
