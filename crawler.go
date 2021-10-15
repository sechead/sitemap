package main

import (
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type crawler struct {
	startPath  string
	host       string
	pages      *Set
	withPrefix string
	stack      *[]string
	withLogs   bool
}

type CrawlerBuilder struct {
	startPath  string
	host       string
	pages      *Set
	withPrefix string
	withLogs   bool
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

func (b *CrawlerBuilder) WithLogs(withLogs bool) *CrawlerBuilder {
	b.withLogs = withLogs
	return b
}

func (b *CrawlerBuilder) Build() *crawler {
	return &crawler{
		startPath:  b.startPath,
		host:       b.host,
		pages:      NewSet(),
		withPrefix: b.withPrefix,
		stack:      &[]string{},
		withLogs:   b.withLogs,
	}
}

func (c crawler) crawl() {
	*c.stack = append(*c.stack, c.startPath)
	for len(*c.stack) > 0 {
		n := len(*c.stack) - 1
		page := (*c.stack)[n]
		*c.stack = (*c.stack)[:n]
		err := c.processPage(&page)
		if err != nil {
			log.Println(err.Error())
			continue
		}
	}
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
		if c.withLogs {
			log.Printf("Found page %s", address)
		}
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
		return nil, errors.New(fmt.Sprintf("status code error for page %s: %d %s\n", *page, res.StatusCode, res.Status))
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
