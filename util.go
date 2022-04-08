package main

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"os"
	"strings"
)

const (
	IllegalCharacters = `\/:*?"<>|`
)

func filename(resp *http.Response) string {
	if dispos := resp.Header.Get("content-disposition"); dispos != "" {
		if _, params, err := mime.ParseMediaType(dispos); err == nil {
			if filename, ok := params["filename"]; ok {
				return removeChars(filename, IllegalCharacters)
			}
		}
	}
	return ""
}

func removeChars(s string, chars string) string {
	for _, ch := range chars {
		s = strings.ReplaceAll(s, string(ch), "")
	}
	return s
}
func urlJoin(base string, pathes ...string) string {
	for _, path := range pathes {
		base = strings.Trim(base, "/")
		path = strings.Trim(path, "/")
		base = base + "/" + path
	}
	return base
}

func wrapperTimeOutError(err error) error {
	switch e := err.(type) {
	case net.Error:
		if e.Timeout() {
			return fmt.Errorf("timeout")
		}
		return err
	default:
		return err
	}
}

func IsValidFilename(fp string) bool {
	// Check if file already exists
	if _, err := os.Stat(fp); err == nil {
		return true
	}

	// Attempt to create it
	var d []byte
	if err := ioutil.WriteFile(fp, d, 0644); err == nil {
		os.Remove(fp) // And delete it
		return true
	}

	return false
}
