package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const (
	logfile = `.dl-download.log`
)

func saveDownloadHistory(tb TaleBook) {
	u, _ := url.Parse(tb.api)

	if err := os.WriteFile(logfile, []byte(fmt.Sprintf("%s %d", u.Host, tb.index-1)), 0644); err != nil {
		log.Printf("warning: can not write dl-download infromat to .dl-download.log %s", err)
	}
}

func tryReadHistoryIndex(api string) (int, error) {
	content, err := os.ReadFile(logfile)
	if err != nil {
		return 0, err
	}
	u, err := url.Parse(api)
	if err != nil {
		return 0, nil
	}
	result := strings.Split(string(content), " ")
	if len(result) >= 2 && result[0] == u.Host {
		return strconv.Atoi(result[1])
	}
	return 0, fmt.Errorf("prase .dl-download.log failed , content: %s", string(content))
}
