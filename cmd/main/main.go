package main

import (
	"log"

	gowikidump "github.com/BehzadE/go-wikidump/pkg/go-wikidump"
)

func main() {
	dump := gowikidump.NewDump()
	dump.Parameters.DumpDirectory = "/home/solaire/Projects/Thesis/dump/"
	err := dump.SetDownloadLinks()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(dump.Links)
	dump.DownloadURLS(3)
}

// func main() {
// err := gowikidump.DownloadURL("https://dumps.wikimedia.orgenwiki-latest-pages-articles-multistream-index1.txt-p1p41242.bz2-rss.xml","/home/solaire/tmp/")
// if err != nil {
//     log.Fatal(err)
// }
// resp, err := grab.Get(".","https://dumps.wikimedia.orgenwiki-latest-pages-articles-multistream-index1.txt-p1p41242.bz2-rss.xml")
// if err != nil {
//     log.Fatal(err)
// }
// fmt.Println("Download saved to",resp.Filename)
// }
