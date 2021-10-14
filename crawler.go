package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type crawler struct {
	startPath  string
	host       string
	pages      *Set
	withPrefix string
	stack      *[]string
}

type CrawlerBuilder struct {
	startPath  string
	host       string
	pages      *Set
	withPrefix string
}

func NewCrawlerBuilder() *CrawlerBuilder {
	return &CrawlerBuilder{}
}

func (b *CrawlerBuilder) WithHost(host string) *CrawlerBuilder {
	b.host = host
	return b
}

func (b *CrawlerBuilder) WithStartPath(startPath string) *CrawlerBuilder {
	b.startPath = startPath
	return b
}

func (b *CrawlerBuilder) WithPrefix(prefix string) *CrawlerBuilder {
	b.withPrefix = prefix
	return b
}

func (b *CrawlerBuilder) Build() *crawler {
	return &crawler{
		startPath:  b.startPath,
		host:       b.host,
		pages:      NewSet(),
		withPrefix: b.withPrefix,
		stack:      &[]string{},
	}
}

func (c crawler) crawl() error {
	*c.stack = append(*c.stack, c.startPath)
	for len(*c.stack) > 0 {
		n := len(*c.stack) - 1
		page := (*c.stack)[n]
		*c.stack = (*c.stack)[:n]
		err := c.processPage(&page)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c crawler) processPage(page *string) error {
	doc, err := c.getPage(page)
	if err != nil {
		return err
	}
	c.processDoc(doc)
	return nil
}

func (c crawler) processDoc(doc *goquery.Document) {
	for _, node := range doc.Find("a").Nodes {
		c.processTag(node)
	}
}

func (c crawler) processTag(node *html.Node) {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			c.processAttribute(attr.Val)
		}
	}
}

func (c crawler) processAttribute(address string) {
	if c.withPrefix != "" {
		if strings.HasPrefix(address, c.withPrefix) {
			c.addPage(address)
		}
	} else {
		c.addPage(address)
	}
}

func (c crawler) addPage(address string) {
	if !c.pages.Has(address) {
		c.pages.Add(address)
		*c.stack = append(*c.stack, address)
	}
}

func (c crawler) getPage(page *string) (*goquery.Document, error) {
	res, err := http.Get(c.host + *page)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Printf("status code error for page %s: %d %s\n", *page, res.StatusCode, res.Status)
		return nil, err
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func crawl() {
	var stack []string
	stack = append(stack, "/")
	var pages = NewSet()
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

	for len(stack) > 0 {
		n := len(stack) - 1
		page := stack[n]
		stack = stack[:n]

		res, err := http.Get("https://sechead.com" + page)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Printf("status code error for page %s: %d %s\n", page, res.StatusCode, res.Status)
			continue
		}

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		for _, node := range doc.Find("a").Nodes {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					if strings.HasPrefix(attr.Val, "/") {
						if !pages.Has(attr.Val) {
							pages.Add(attr.Val)
							_, err := w.WriteString(
								fmt.Sprintf(
									"<url><loc>https://sechead.com%s</loc></url>\n",
									strings.ReplaceAll(attr.Val, "&", "%26"),
								),
							)
							if err != nil {
								log.Fatal(err)
							}
							fmt.Println(attr.Val)
							stack = append(stack, attr.Val)
						}
					}
				}
			}
		}
	}
	_, err = w.WriteString("</urlset>")
	if err != nil {
		log.Fatal(err)
	}
	err = w.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
