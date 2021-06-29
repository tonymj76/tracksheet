package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const URL = "https://www.amazon.com/b?node=16225016011&pf_rd_r=7HGEC8DG3WZ0MMT709BQ&pf_rd_p=e5b0c85f-569c-4c90-a58f-0c0a260e45a0&pd_rd_r=5f38106e-4215-41f6-b427-d75f52194c31&pd_rd_w=T9rga&pd_rd_wg=EysZ8&ref_=pd_gw_unk"

func ExampleScrape() {
	// Request the HTML page.
	res, err := http.Get(URL)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	urls := []string{}
	// Find the review items
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		// title := s.Find("a").Text()
		t, _ := s.Attr("href")
		// fmt.Println(s.Attr("href"))
		urls = append(urls, t)
	})
	fmt.Println(len(urls))
	result := GetUnique(urls, "/")
	fmt.Println(result)
	fmt.Println(len(result))
}

func main() {
	ExampleScrape()
}
