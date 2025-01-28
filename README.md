# Go-webcrawler

Here’s a basic README for your web crawler project:

Web Crawler:-

This is a simple web crawler written in Go that extracts the titles and URLs from a website and saves the data in a CSV file.
The crawler is visiting web pages concurrently to increase crawling efficiency .

Features:-
- Crawls a website starting from a given URL.
- Extracts page titles and saves the URL and title into a CSV file.
- Multiple pages are crawled at the same time to speed up the crawling process.
- The program uses Go channels to synchronize concurrent crawling tasks, ensuring all pages are crawled before the program exits.
- Implements rate limiting with a delay between requests to avoid overwhelming the server.

Setup:-
1. Clone the repository:
    git clone <repository-url>

2. Navigate to the project directory:
    cd web-crawler
   
4. Install dependencies:
    go get github.com/gocolly/colly

5. Run the web crawler:
    go run webcrawler.go

Output:-
The crawler generates a "crawled_data.csv" file with the following columns:
->URL – The URL of the crawled page.
->Title – The title of the crawled page.
