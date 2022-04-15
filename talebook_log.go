package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

const (
	name = `.dl-tabebook.json`
)

func saveDownloadHistory(tb TaleBook) {
	home, _ := os.UserHomeDir()
	logfile := filepath.Join(home, name)
	u, _ := url.Parse(tb.api)

	data, err := readjsonMap(logfile)
	if err != nil {
		log.Printf("history not foud,create new history now")
		data = make(map[string]int)
	}
	data[u.Host] = tb.index - 1
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("warn gen history json failed")
	}
	if err = os.WriteFile(logfile, content, 0644); err != nil {
		log.Printf("warning: can not write dl-download infromat to .dl-download.log %s", err)
	}
}

func readjsonMap(filename string) (map[string]int, error) {
	home, _ := os.UserHomeDir()
	logfile := filepath.Join(home, name)
	content, err := os.ReadFile(logfile)
	if err != nil {
		return nil, err
	}
	var data = make(map[string]int)
	err = json.Unmarshal(content, &data)
	return data, err
}
func tryReadHistoryIndex(api string) (int, error) {
	home, _ := os.UserHomeDir()
	logfile := filepath.Join(home, name)
	data, err := readjsonMap(logfile)
	if err != nil {
		return 0, err
	}
	u, err := url.Parse(api)
	if err != nil {
		return 0, nil
	}
	if index, ok := data[u.Host]; ok {
		return index, nil
	}
	return 0, fmt.Errorf("download history not found")
}
