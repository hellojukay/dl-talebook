package main

import (
	"flag"
	"log"
	"time"
)

var (
	userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36"
	site      = `https://book.codefine.site:6870/`
	dir       = "./"
	timeout   = time.Duration(10) * time.Second
	username  = ""
	password  = ""
)

func init() {
	flag.StringVar(&username, "username", username, "username")
	flag.StringVar(&password, "password", password, "password")
	flag.StringVar(&site, "site", site, "tabebook web site")
	flag.StringVar(&dir, "dir", dir, "data dir")
	flag.DurationVar(&timeout, "timeout", timeout, "http timeout")
	flag.StringVar(&userAgent, "user-agent", userAgent, "http userAgent")

	flag.Parse()
}
func main() {
	tale, err := NewTableBook(site, WithTimeOutOption(timeout), WithLoginOption(username, password))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d books retrieved on server %s", tale.ServerInfo.Sys.Books, site)
	for {
		book, err := tale.Next()
		if err != nil {
			log.Printf("%s %s [skiped]", site, err.Error())
			continue
		}
		log.Printf("downloading %s", book.String())
		if err = tale.Download(book, dir); err != nil {
			log.Printf("%s %s [skiped]", book.Book.Title, err)
		}
	}
}
