package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	crawler := NewCrawlerBuilder().
		WithHost("https://sechead.com").
		WithStartPath("/").
		WithPrefix("/").
		Build()
	err := crawler.crawl()
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("sitemap.xml")
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(f)
	defer f.Close()
	_, err = w.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?><urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")
	if err != nil {
		log.Fatal(err)
	}
	crawler.pages.Each(func(item string) bool {
		_, err := w.WriteString(item)
		if err != nil {
			log.Fatal(err)
		}
		return true
	})
	_, err = w.WriteString("</urlset>")
	if err != nil {
		log.Fatal(err)
	}
	err = w.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
