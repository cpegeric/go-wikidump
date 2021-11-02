package download

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cavaliercoder/grab"
)

func getLink(link string, dst string) error {
	client := grab.NewClient()
	req, err := grab.NewRequest(dst, link)
	if err != nil {
		return err
	}
	urlSplit := strings.Split(link, "/")
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

func worker(urls <-chan string, results map[string]error, dst string, wg *sync.WaitGroup) {
	for url := range urls {
		err := getLink(url, dst)
		results[url] = err
	}
	wg.Done()
}
