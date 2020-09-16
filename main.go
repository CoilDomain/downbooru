package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cavaliercoder/grab"

	"github.com/PuerkitoBio/goquery"
)

// Create function to grab content of page and image URLs
//              Input ;)     Output ;)
func datascrape(URL string) {
	// Getting the content of the page
	res, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	// Searching for the images URLs
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".post-preview").Each(func(i int, s *goquery.Selection) {
		URL, _ := s.Attr("data-file-url")
		// Input scraped URL into database
		dbinsert(URL)
	})
}

// Set download location
var foldername = "Downloads"
var downloadfolder = filepath.Join(path, foldername)

// Download inage function, need to make semi-async within limits
func getimg(URL string) {
	resp, err := grab.Get(downloadfolder, URL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error downloading ", URL, err)
		os.Exit(1)
	}
	fmt.Println("Downloading: ", resp.Filename)
}

// Poll database and download images, possibly limit first 10, mark complete and reiterate for faster downloads
func poll() {
	for row := range dbquery() {
		getimg(row)
	}
}

// Main function
func main() {
	// Commandline arguments
	tagPtr := flag.String("tags", "robot", "Any Danbooru tag")
	pagesPtr := flag.String("pages", "1", "Number of pages to go through")
	flag.Parse()

	// Creating the proper URL to scrape
	baseURL := "https://danbooru.donmai.us/posts?"

	// Create variable to split the tags into separate strings to use in a foreach statement
	tags := strings.Split(*tagPtr, ",")
	// Create loop for multiple tags
	for _, tag := range tags {
		// Convert the page argument variable into an into so I can do a numeric compare
		p, _ := strconv.Atoi(*pagesPtr)
		// Run loop on pages to back track from the selected page
		for n := 1; n < p+1; n++ {
			// Creating the proper URL to scrape
			pns := strconv.Itoa(n)
			pageURL := "page=" + pns
			tagURL := "&tags=" + tag
			fullURL := baseURL + pageURL + tagURL
			datascrape(fullURL)
		}
	}
	poll()
}
