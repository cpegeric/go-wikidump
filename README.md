# go-wikidump
This package was written in order to facilitate working with mediawiki dumps. You can use this package
to download a mediawiki dump and extract wikipedia pages from bz2 compressed files without the need to
extract the files (which can be quite big). 
## Installation 
    go get https://github.com/BehzadE/go-wikidump

## Usage
Initialize a dump struct:

    dump := gowikidump.NewDump()
This initializes a dump with the default parameters. You can change the default paramters.

    dump.Parameters.DumpVer = "/enwiki/20210901/"
If you have already downloaded pages-articles-multistream dump you can set the DumpDirectory
to point to the download directory. Otherwise set the download links in the dump struct:

    err := dump.SetDownloadLinks()
and download the dump:

    err =  dump.DownloadURLS(3)
where the number specifies the number of concurrent downloads. According to my tests trying
to use more than 3 concurrent download connections results in errors so choose from 1-3.

The default download folder is "wikipedia-dump" inside the working directory.
Before you can start extracting pages you need to run 
    
    dump.SaveIndexRanges()
Which reads the index files and saves the range of pageIDs included in each index in a text file
to be able to access the required index file without having to search through all of them for 
each page. This function only needs to be called once to create offsets.txt file in the dump directory.


Extract a specific page from the dump:

    pageID = 12
    page, err := dump.GetPage(pageID)
On unix systems with pandoc installed you can extract the plain text article using:

    plain, err := page.GetPlainText()
