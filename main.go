package main

import (
	"flag"
	"log"
	"time"
)

var (
	site    = `https://book.codefine.site:6870/`
	dir     = "./"
	timeout = time.Duration(10) * time.Second
)

func init() {
	flag.StringVar(&site, "site", site, "tabebook web site")
	flag.StringVar(&dir, "dir", dir, "data dir")
	flag.DurationVar(&timeout, "timeout", timeout, "http timeout")
	flag.Parse()
}
func main() {
	tale, err := NewTableBook(site, WithTimeOutOption(timeout))
	if err != nil {
		log.Fatal(err)
	}
	for {
		book, err := tale.Next()
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Printf("downloading %s", book.String())
		if err = tale.Download(book, dir); err != nil {
			log.Fatalf("download %s , %s", book.Book.Title, err)
		}
	}
}
