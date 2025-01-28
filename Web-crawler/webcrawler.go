package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

var visited = make(map[string]bool)

func crawl(currenturl string, depth int) {
	//CSV file to store data
	file, err := os.Create("crawled_data.csv")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Creating a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"URL", "Title"})

	c := colly.NewCollector(
		colly.AllowedDomains("www.scrapingcourse.com", "scrapingcourse.com"),
		colly.MaxDepth(depth),
	)

	// Fetching title and writing to the CSV
	c.OnHTML("title", func(e *colly.HTMLElement) {
		title := e.Text
		url := e.Request.URL.String()
		fmt.Println("Title:", title)

		// Writing data to the CSV file
		writer.Write([]string{url, title})
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link != "" && !visited[link] && (len(link) > 4 && link[len(link)-4:] != ".pdf" && link[len(link)-4:] != ".png") {
			visited[link] = true
			fmt.Println("Link found:", link)
			e.Request.Visit(link)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Crawling website:", r.URL)
	})

	err = c.Visit(currenturl)
	if err != nil {
		fmt.Print("Error while visiting:", err)
	}
}

func main() {
	seedurl := "https://www.scrapingcourse.com/ecommerce/"
	crawl(seedurl, 3)
	fmt.Println("Crawling Summary:")
	fmt.Println("No of pages visisted : ", len(visited))
}
