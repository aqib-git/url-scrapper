package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Scrapper struct {
	depth int;
	urls map[string]int;
}

func (scrapper Scrapper) Scrap (url string, depth int) {
	if depth > scrapper.depth || !strings.HasPrefix(url, "http") {
		return
	}

	if _, found := scrapper.urls[url]; found && scrapper.urls[url] > 0 {
		return
	}

	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)

		return
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
	  log.Fatal(err)

	  return
	}

	urls := []string{}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		url, ok := s.Attr("href")

		if !ok {
			return
		}


		if _, found := scrapper.urls[url]; found {
			scrapper.urls[url] += 1

			return
		}

		if !strings.HasPrefix(url, "http") {
			return
		}

		urls = append(urls, url)

		scrapper.urls[url] = 0
	})

	for i := 0; i < len(urls); i++ {
		scrapper.Scrap(urls[i], depth + 1)
	}
}

func main() {
	scrapy := Scrapper{depth: 2, urls: map[string]int{}}
	scrapy.Scrap("http://aqibpandit.com", 1)
	
	urlBytes, err := json.Marshal(scrapy.urls)

	if err != nil {
		log.Fatal(err)
		
		return
	}

	fmt.Println(string(urlBytes))
}
