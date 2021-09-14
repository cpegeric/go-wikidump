package gowikidump

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func SavePageOffsets(dumpDirectory string) error {
	indexFiles, err := GetIndexFiles(dumpDirectory)
	indexDicts, err := GetPageIDRanges(indexFiles)
	if err != nil {
		return err
	}
	file, err := os.Create(dumpDirectory + "/offsets.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	sort.Slice(indexDicts, func(i int, j int) bool {
		int1, err := strconv.Atoi(indexDicts[i].EndID)
		if err != nil {
			panic(err)
		}
		int2, err := strconv.Atoi(indexDicts[j].EndID)
		if err != nil {
			panic(err)
		}
		return int1 < int2
	})

	for _, row := range indexDicts {
		fmt.Fprintln(file, row.EndID+"###"+row.Filename)
	}
	return nil
}

func GetPageIDRanges(indexFiles []string) ([]Index, error) {
	indexDicts := make([]Index, 0)
	for _, indexFile := range indexFiles {
		tailCmd := "bzcat " + indexFile + " | tail -1"
		tail, err := exec.Command("sh", "-c", tailCmd).Output()
		if err != nil {
			return nil, err
		}
		tailID := strings.Split(string(tail), ":")[1]

		if err != nil {
			return nil, err
		}
		toAdd := Index{
			EndID:    tailID,
			Filename: indexFile,
		}
		indexDicts = append(indexDicts, toAdd)
	}
	return indexDicts, nil
}

func GetIndexFiles(dumpDirectory string) ([]string, error) {
	cmd := exec.Command("find", dumpDirectory, "-name", "*index*")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	result := strings.Split(out.String(), "\n")
	return result[:len(result)-1], nil
}

func GetPageBZ2File(dumpDirectory string, pageID int64) (string, error) {
	file, err := os.Open(dumpDirectory + "/offsets.txt")
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

func GetPageByteLocation(indexFilename string, pageID int64) ([]int64, error) {
	cmd := "bzcat " + indexFilename + " | rg :" + fmt.Sprint(pageID) + ":" + " --line-number"
	line, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return nil, err
	}
	linesplit := strings.Split(string(line), ":")
	linenum, err := strconv.Atoi(linesplit[0])
	if err != nil {
		return nil, err
	}
	byteBegin, err := strconv.ParseInt(linesplit[1], 10, 64)
	if err != nil {
		return nil, err
	}
	cmd = "bzcat " + indexFilename + " | sed '" + fmt.Sprint(linenum+100) + "q;d'"
	line, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return nil, err
	}
	linesplit = strings.Split(string(line), ":")
	byteEnd, err := strconv.ParseInt(linesplit[0], 10, 64)
	if err != nil {
		return nil, err
	}

	return []int64{byteBegin, byteEnd}, nil
}

func GetPageStream(indexFilename string, byteLocations []int64) ([]byte, error) {
	dataFilename := strings.Replace(indexFilename, "txt", "xml", 1)
	dataFilename = strings.Replace(dataFilename, "-index", "", 1)
	file, err := os.Open(dataFilename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	sr := io.NewSectionReader(file, byteLocations[0], byteLocations[1]-byteLocations[0])
	reader := bzip2.NewReader(sr)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetPageFromStream(stream []byte, pageID int64) (*Page, error) {
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

func GetPage(dumpLocation string, pageID int64) (*Page, error) {
	filename, err := GetPageBZ2File(dumpLocation, pageID)
	if err != nil {
		return nil, err
	}
	byteLocations, err := GetPageByteLocation(filename, pageID)
	if err != nil {
		return nil, err
	}

	stream, err := GetPageStream(filename, byteLocations)
	if err != nil {
		return nil, err
	}

	page, err := GetPageFromStream(stream, int64(pageID))
	if err != nil {
		return nil, err
	}
	return page, nil
}

func GetPages(dumpLocation string, pageIDs []int64) ([]*Page, error) {
	result := make([]*Page, len(pageIDs))
	for index, id := range pageIDs {
		page, err := GetPage(dumpLocation, id)
		if err != nil {
			return nil, err
		}
		result[index] = page
	}
	return result, nil
}

func (page *Page) GetSectionTitles() []string {
	r := regexp.MustCompile("(?m)^== .*? ==$")
	matches := r.FindAllString(page.Revision.Text, -1)
	return matches
}

func NormalizeSectionTitle(title string) string {
	r := regexp.MustCompile("^== (.*?) ==$")
	rtrim := regexp.MustCompile("^[a-zA-Z0-9 ,]*")
	return rtrim.FindString(r.FindStringSubmatch(title)[1])
}

func (page *Page) GetPlainText() ([]byte, error) {
	cmd := exec.Command("sh", "-c", "pandoc -f mediawiki -t plain")
	bin := bytes.Buffer{}
	bin.Write([]byte(page.Revision.Text))
	bout := bytes.Buffer{}
	cmd.Stdin = &bin
	cmd.Stdout = &bout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return bout.Bytes(), nil
}

func Set(ptr interface{}, tag string) error {
	if reflect.TypeOf(ptr).Kind() != reflect.Ptr {
		return fmt.Errorf("Not a pointer")
	}
	v := reflect.ValueOf(ptr).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		if defaultVal := t.Field(i).Tag.Get(tag); defaultVal != "-" {
			if err := setField(v.Field(i), defaultVal); err != nil {
				return err
			}
		}
	}
	return nil
}

func setField(field reflect.Value, defaultVal string) error {

	if !field.CanSet() {
		return fmt.Errorf("Can't set value\n")
	}
	switch field.Kind() {
	case reflect.Int:
		if val, err := strconv.ParseInt(defaultVal, 10, 64); err == nil {
			field.Set(reflect.ValueOf(int(val)).Convert(field.Type()))
		}
	case reflect.String:
		field.Set(reflect.ValueOf(defaultVal).Convert(field.Type()))
	}
	return nil
}
