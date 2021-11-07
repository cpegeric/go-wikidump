# go-wikidump
## Introduction
Wikipedia dumps are a great source of textual data for language processing and machine learning purposes. The aim of this module is to make extracting pages from such dumps easier. This module works with the multistream xml dumps. Refer to [Wikipedia](https://en.wikipedia.org/wiki/Wikipedia:Database_download#Should_I_get_multistream?) for more information on the multistream dumps.

In short, multstream dumps are made up of streams, each one of which holds only 100 wikipedia pages. Each multistream dump file comes with an index file of the same name. Index files contain the byte locations of the stream for pages. 
## Features
- Save the index file information into a sqlite database for easier quering of data.
- Extract individual streams using the byte locations without extracting the whole file.
- Parse the xml in stream to get individual pages
## TODO 
- Writing Tests
- Optimization
- Removing or expanding the templates in wikitext.
- Parsing wikitext to plain text.
- Parsing wikitext to html.
## Installation 
    go get https://github.com/BehzadE/go-wikidump

## Usage
Download one or all parts of the multistream wikipedia dump into a directory. Each file must come with the corresponding index file.

    import (
        "fmt"
        "log"

        "github.com/BehzadE/go-wikidump/pkg/wikidump"
    )

    func main() {
        path := "/home/solaire/Data/wikidump/"
        d, err := wikidump.New(path)
        if err != nil {
            log.Fatal(err)
        }
        err = d.PopulateDB()
        if err != nil {
            log.Fatal(err)
        }
        pages, err := d.GetPages([]int64{12, 13, 14, 15, 622, 624, 1941, 1944})
        if err != nil {
            log.Fatal(err)
        }
        for _, page := range pages {
            fmt.Println(page.Revision.Text)
        }
    }

This will find the pages with the given IDs. PopulateDB only needs to be called once and has no effect if called again.

You can also get a streamReader for a given dump file to read streams one by one from the begining of the file:

	reader, err := d.NewStreamReader("enwiki-20210720-pages-articles-multistream1.xml-p1p41242.bz2")
	if err != nil {
		log.Fatal(err)
	}
	for reader.Next() {
		b, err := reader.Read()
		if err != nil {
			log.Fatal(err)
		}
		pages, err := wikidump.ParseStream(b)
		if err != nil {
			log.Fatal(err)
		}
		for _, page := range pages {
			fmt.Println(page.Title)
		}
	
	}

Although you'd be better off extracting the file and parsing the xml instead of going through all the extra steps to extract individual streams if you plan on reading the whole file.
