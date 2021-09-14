package gowikidump

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cavaliercoder/grab"
	"golang.org/x/net/html"
)

func NewDump() *Dump {
	dump := Dump{}
	Set(&dump.Parameters, "default")
	return &dump
}

func (dump *Dump) SetDownloadLinks() error {
	dumpURL := dump.Parameters.BaseURL + dump.Parameters.DumpVer
	resp, err := http.Get(dumpURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	links := getHtmlLinks(resp.Body)[2:]
	for i := range links {
		links[i] = dump.Parameters.BaseURL + links[i]
	}
	dump.Links = links
	return nil
}

func getHtmlLinks(body io.Reader) []string {
	var links []string
	z := html.NewTokenizer(body)
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.EndTagToken:
			token := z.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						if strings.Contains(attr.Val, "pages-articles-multistream") {
							links = append(links, attr.Val)
						}
					}
				}
			}
		}
	}
}

func (dump *Dump) DownloadURLS(maxWorkers int) error {
	if dump.Links == nil {
		return errors.New("Links unset.")
	}
	err := os.MkdirAll(dump.Parameters.DumpDirectory, os.ModePerm)

	if err != nil {
		return err
	}

	results := make(map[string]error)
	linksC := make(chan string, len(dump.Links))
	var wg sync.WaitGroup
	for _, link := range dump.Links {
		linksC <- link
	}
	close(linksC)
	for w := 1; w <= maxWorkers; w++ {
		wg.Add(1)
		go DownloadWorker(w, linksC, results, dump.Parameters.DumpDirectory, &wg)
	}
	wg.Wait()
	for k, v := range results {
		if v != nil {
			fmt.Printf("Url %v failed with error: %v.\n", k, v)
		}
	}
	return nil
}

func DownloadURL(url string, dst string) error {
	client := grab.NewClient()
	req, err := grab.NewRequest(dst, url)
	if err != nil {
		return err
	}
	urlSplit := strings.Split(url, "/")
	filename := urlSplit[len(urlSplit)-1]
	fmt.Printf("Downloading %v...", filename)
	resp := client.Do(req)
	fmt.Println(resp.HTTPResponse.Status)
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
Loop:
	for {
		select {
		case <-t.C:
			fmt.Println(filename)
			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress())

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	if err := resp.Err(); err != nil {
		return err
	}
	return nil
}

func DownloadWorker(id int, urls <-chan string, results map[string]error, dst string, wg *sync.WaitGroup) {
	for url := range urls {
		err := DownloadURL(url, dst)
		results[url] = err
	}
	wg.Done()
}
