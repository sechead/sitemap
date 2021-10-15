package main

import (
	"github.com/robfig/cron/v3"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	go sitemapUpdater()
	fs := http.FileServer(http.Dir("./data"))
	http.Handle("/", fs)
	log.Printf("Listening on :%s\n", os.Getenv("PORT"))
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal(err.Error())
	}

}

func sitemapUpdater() {
	updateSitemap()
	c := cron.New(cron.WithLogger(cron.VerbosePrintfLogger(log.Default())))
	_, err := c.AddFunc(os.Getenv("SCHEDULE"), func() { updateSitemap() })
	if err != nil {
		log.Fatal(err.Error())
	}
	c.Start()
}

func updateSitemap() {
	withLogs, err := strconv.ParseBool(os.Getenv("WITH_LOGS"))
	if err != nil {
		log.Fatal(err)
	}
	crawler := NewCrawlerBuilder().
		WithHost(os.Getenv("HOST")).
		WithStartPath(os.Getenv("START_PATH")).
		WithPrefix(os.Getenv("PREFIX")).
		WithLogs(withLogs).
		Build()
	crawler.crawl()
	if err != nil {
		log.Fatal(err)
	}
	newpath := filepath.Join(".", "data")
	err = os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	err = writeToFile(crawler.pages, os.Getenv("HOST"), "data/sitemap_new.xml")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove("data/sitemap.xml")
	if err != nil {
		log.Println(err)
	}
	err = os.Rename("data/sitemap_new.xml", "data/sitemap.xml")
	if err != nil {
		log.Fatal(err)
	}
}
