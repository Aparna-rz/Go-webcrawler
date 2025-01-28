package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

var visited = make(map[string]bool)
var wg sync.WaitGroup // WaitGroup to wait for all goroutines to finish

func crawl(currenturl string, depth int, writer *csv.Writer) {
	defer wg.Done() // Decrement the counter when the function completes

	c := colly.NewCollector(
		colly.AllowedDomains("www.scrapingcourse.com"),
		colly.MaxDepth(depth),
	)

	// Apply rate limiting (max 1 request per second)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       1 * time.Second,        // Delay of 1 second between requests
		RandomDelay: 500 * time.Millisecond, // Random delay to avoid hitting server too predictably
	})

	// Fetching title and writing to the CSV
	c.OnHTML("title", func(e *colly.HTMLElement) {
		title := e.Text
		url := e.Request.URL.String()
		fmt.Println("Title:", title)

		// Write data to the CSV file
		writer.Write([]string{url, title})
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link != "" && !visited[link] && (len(link) > 4 && link[len(link)-4:] != ".pdf" && link[len(link)-4:] != ".png") {
			visited[link] = true
			fmt.Println("Link found:", link)

			// Visit the new link concurrently
			wg.Add(1)                     // Increment WaitGroup counter for the new goroutine
			go crawl(link, depth, writer) // Launch a new goroutine for the link
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Crawling website:", r.URL)
	})

	// Start visiting the URL
	err := c.Visit(currenturl)
	if err != nil {
		fmt.Print("Error while visiting:", err)
	}
}

func main() {
	// Open a CSV file to store data (created once)
	file, err := os.Create("crawled_data.csv")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header row
	writer.Write([]string{"URL", "Title"})

	// Start crawling
	seedurl := "https://www.scrapingcourse.com/ecommerce/"
	wg.Add(1)                    // Add 1 to WaitGroup counter for the initial URL
	go crawl(seedurl, 1, writer) // Start crawling the seed URL concurrently

	// Wait for all goroutines to finish
	wg.Wait()

	// Print crawling summary
	fmt.Println("Crawling Summary:")
	fmt.Println("No of pages visited: ", len(visited))
}
