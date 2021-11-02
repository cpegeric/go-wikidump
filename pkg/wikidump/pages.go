package wikidump

import (
	"bytes"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// Get the plain text from the original wikitext. Requires the pandoc package to be installed on the system.
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

func (page *Page) getSectionTitleLines() []string {
	r := regexp.MustCompile("(?m)^=+.*?=+$")
	matches := r.FindAllString(page.Revision.Text, -1)
	return matches
}

// Get page sections.
func (page *Page) GetSectionTitles() []*Section {
	lines := page.getSectionTitleLines()
	sections := make([]*Section, 0)
	for _, line := range lines {
		level := strings.Count(line, "=") / 2
		r := regexp.MustCompile("^=+(.*?)=+$")
		rtrim := regexp.MustCompile("^[a-zA-Z0-9 ,]*")
		title := rtrim.FindString(r.FindStringSubmatch(line)[1])
		sections = append(sections, &Section{
			Title: strings.TrimSpace(title),
			Level: level,
		})
	}
	return sections
}
