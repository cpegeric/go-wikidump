package gowikidump

import (
	"bytes"
	"os"
	"os/exec"
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
