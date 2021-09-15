package main

import (
	"fmt"
	"log"

	gowikidump "github.com/BehzadE/go-wikidump/pkg/go-wikidump"
)

func main() {
	dump := gowikidump.NewDump()
	dump.Parameters.DumpDirectory = "/home/solaire/Projects/Thesis/dump/"
	// err := dump.SetDownloadLinks()
	// dump.DownloadURLS(3)
	// err := dump.SaveIndexRanges()
	var pageID int64
	pageID = 57027716
	// pageID = 12
	page, err := dump.GetPage(pageID)
	if err != nil {
		log.Fatal(err)
	}
	plain, err := page.GetPlainText()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(plain))

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
