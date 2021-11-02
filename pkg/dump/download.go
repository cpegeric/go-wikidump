package dump

import (
	"io"
	"net/http"

	"github.com/BehzadE/go-wikidump/pkg/download"
)

// Find and save the download urls for a given version of mediawiki dump.
func (dump *dump) Download(maxWorkers int) error {
	dumpURL := dump.BaseURL + dump.Ver
	resp, err := http.Get(dumpURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	links := getMultistream(resp.Body)
	for i := range links {
		links[i] = dump.BaseURL + links[i]
	}
	err = download.GetLinks(links, dump.Directory, maxWorkers)
	return err
}

func getMultistream(body io.Reader) []string {
	return download.ExtractLinks(body, "pages-articles-multistream")[2:]
}
