package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"mime"
	"net"
	"net/http"
	"os"
	"strings"
)

func filename(resp *http.Response) string {
	const (
		IllegalCharacters = `\/:*?"<>|`
	)

	if dispos := resp.Header.Get("content-disposition"); dispos != "" {
		if _, params, err := mime.ParseMediaType(dispos); err == nil {
			if filename, ok := params["filename"]; ok {

				// https://github.com/hellojukay/dl-talebook/issues/5
				return removeChars(filename, IllegalCharacters)
			}
		}
	}
	return ""
}

func removeChars(s string, chars string) string {
	for _, ch := range chars {
		s = strings.ReplaceAll(s, string(ch), " ")
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

func Bytes(s uint64) string {
	sizes := []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(s, 1000, sizes)
}

func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f %s"
	if val < 10 {
		f = "%.1f %s"
	}

	return fmt.Sprintf(f, val, suffix)
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}
