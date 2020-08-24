package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	// Commandline arguments
	tagPtr := flag.String("tag", "robot", "Any Danbooru tag")
	pagesPtr := flag.String("pages", "1", "Number of pages to go through")
	flag.Parse()
	// Creating the proper URL to scrape
	baseURL := "https://danbooru.donmai.us/posts?"
	pages := "page=" + *pagesPtr
	tags := "&tags=" + *tagPtr
	fullURL := baseURL + pages + tags
	// Creating a variable to use as the page int
	p, _ := strconv.Atoi(*pagesPtr)
	// Check whether p is larger than 1 so I can do a sequence for pages before p
	if p > 1 {
		// Do a loop based on how many pages are between 1 and the actual argument
		for n := 1; n < p+1; n++ {
			// Creating the proper URL to scrape
			pns := strconv.Itoa(n)
			pages := "page=" + pns
			fullURL := baseURL + pages + tags
			// Getting the content of the page
			res, err := http.Get(fullURL)
			if err != nil {
				log.Fatal(err)
			}
			defer res.Body.Close()
			if res.StatusCode != 200 {
				log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
			}
			// Searching for the images
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			doc.Find(".post-preview").Each(func(i int, s *goquery.Selection) {
				URL, _ := s.Attr("data-file-url")
				fmt.Printf("%s\n", URL)
			})
		}
	} else {
		// Getting the content of the page
		res, err := http.Get(fullURL)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}
		// Searching for the images
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		doc.Find(".post-preview").Each(func(i int, s *goquery.Selection) {
			URL, _ := s.Attr("data-file-url")
			fmt.Printf("%s\n", URL)
		})
	}

}
