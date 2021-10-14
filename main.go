package main

import (
	"log"
)

func main() {
	const host = "https://sechead.com"
	crawler := NewCrawlerBuilder().
		WithHost(host).
		WithStartPath("/").
		WithPrefix("/").
		WithLogs(true).
		Build()
	err := crawler.crawl()
	if err != nil {
		log.Fatal(err)
	}
	err = writeToFile(crawler.pages,  host, "sitemap.xml")
	if err != nil {
		log.Fatal(err)
	}
}
