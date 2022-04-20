package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"time"
)

var (
	userAgent       = flag.String("user-agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36`, "http user-agent")
	site            = flag.String("site", `https://book.codefine.site:6870/`, "tabebook web site")
	cookie          = flag.String("cookie", "", "http cookie")
	dir             = flag.String("dir", "./", "data dir")
	timeout         = flag.Duration("timeout", time.Duration(10)*time.Second, "http timeout")
	username        = flag.String("username", "", "username")
	password        = flag.String("password", "", "password")
	verbose         = flag.Bool("verbose", false, "show debug log")
	startIndex      = flag.Int("start-index", 0, "start book id")
	version         = flag.Bool("version", false, "show progream version")
	continueOnStart = flag.Bool("continue", true, "continue an incomplete download")
	retry           = flag.Int("retry", 3, "timeout retries count")
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	if !flag.Parsed() {
		flag.Parse()
	}

	if *version {
		PrintVersion()
		os.Exit(0)
	}
}

func main() {
	tale, err := NewTableBook(*site,
		WithRetry(*retry),
		WithVerboseOption(*verbose),
		WithUserCookieOption(*cookie),
		WithUserAgentOption(*userAgent),
		WithTimeOutOption(*timeout),
		WithStartIndex(*startIndex),
		WithLoginOption(*username, *password),
		WithContinue(*continueOnStart),
	)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d books retrieved on server %s", tale.Total, *site)

	for {
		book, err := tale.Next()
		if err != nil {
			if errors.Is(err, NO_MORE_BOOK_ERROR) {
				log.Printf("%s [exit]", err.Error())
				if tale.exit != nil {
					tale.exit()
					os.Exit(0)
				}
			}
			log.Printf("%s [skiped]", err.Error())
			continue
		}

		if err = tale.Download(book, *dir); err != nil {
			log.Printf("[%d/%d] downloading %s, %s [skiped]", book.Book.ID, tale.LastIndex(), book.Book.Title, err)
			continue
		}

		log.Printf("[%d/%d] downloading %s successed", book.Book.ID, tale.LastIndex(), book.String())
	}
}
