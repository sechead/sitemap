package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func writeToFile(set *Set, host string, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	defer f.Close()
	_, err = w.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n")
	if err != nil {
		log.Fatal(err)
	}
	set.Each(func(item string) bool {
		_, err := w.WriteString(
			fmt.Sprintf(
				"<url><loc>%s%s</loc></url>\n",
				host,
				strings.ReplaceAll(item, "&", "%26"),
				),
			)
		if err != nil {
			log.Fatal(err)
		}
		return true
	})
	_, err = w.WriteString("</urlset>")
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	return nil
}
