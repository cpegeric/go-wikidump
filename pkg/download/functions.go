package download

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

// Extract links containing the given string in a html body.
func ExtractLinks(body io.Reader, mustContain string) []string {
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
						if strings.Contains(attr.Val, mustContain) {
							links = append(links, attr.Val)
						}
					}
				}
			}
		}
	}
}

func GetLinks(links []string, dst string, maxWorkers int) error {
	err := os.MkdirAll(dst, os.ModePerm)
	if err != nil {
		return err
	}

	results := make(map[string]error)
	c := make(chan string, len(links))
	var wg sync.WaitGroup
	for _, link := range links {
		c <- link
	}
	close(c)
	for w := 1; w <= maxWorkers; w++ {
		wg.Add(1)
		go worker(c, results, dst, &wg)
	}
	wg.Wait()
	for k, v := range results {
		if v != nil {
			fmt.Printf("Url %v failed with error: %v.\n", k, v)
		}
	}
	return nil
}
