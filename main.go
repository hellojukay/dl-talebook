package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"time"

	"golang.org/x/time/rate"
)

var (
	userAgent  = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36"
	site       = `https://book.codefine.site:6870/`
	dir        = "./"
	timeout    = time.Duration(10) * time.Second
	username   = ""
	password   = ""
	concurrent = 5
)

func init() {
	flag.IntVar(&concurrent, "c", concurrent, "maximum number of concurrent download tasks allowed per second")
	flag.StringVar(&username, "username", username, "username")
	flag.StringVar(&password, "password", password, "password")
	flag.StringVar(&site, "site", site, "tabebook web site")
	flag.StringVar(&dir, "dir", dir, "data dir")
	flag.DurationVar(&timeout, "timeout", timeout, "http timeout")
	flag.StringVar(&userAgent, "user-agent", userAgent, "http userAgent")

	flag.Parse()
}
func main() {
	tale, err := NewTableBook(site,
		WithTimeOutOption(timeout),
		WithLoginOption(username, password),
		WithUserAgentOption(userAgent),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d books retrieved on server %s", tale.ServerInfo.Sys.Books, site)
	l := rate.NewLimiter(rate.Limit(concurrent), concurrent)
	for {
		// 限制速度
		l.Wait(context.Background())
		book, err := tale.Next()
		if err != nil {
			log.Printf("%s [skiped]", err.Error())
			if errors.Is(err, NO_MORE_BOOK_ERROR) {
				os.Exit(0)
			}
			continue
		}
		log.Printf("downloading %s", book.String())
		go func() {
			if err = tale.Download(book, dir); err != nil {
				log.Printf("%s %s [skiped]", book.Book.Title, err)
			}
		}()
	}
}
