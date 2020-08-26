package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Create function to grab content of page and image URLs
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
		fmt.Printf("%s\n", URL)
	})
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
}
