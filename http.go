package main

import (
	"fmt"
	"mime"
	"net"
	"net/http"
	"strings"
)

func filename(resp *http.Response) string {
	if dispos := resp.Header.Get("content-disposition"); dispos != "" {
		if _, params, err := mime.ParseMediaType(dispos); err == nil {
			if filename, ok := params["filename"]; ok {
				return filename
			}
		}
	}
	return ""
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
