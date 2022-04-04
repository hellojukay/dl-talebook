package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

type TaleBook struct {
	api    string
	index  int
	client *http.Client
	err    error
}

type Book struct {
	ERROR string `json:"err"`
	Book  struct {
		Title string `json:"title"`
		FILES []struct {
			Format string `json:"format"`
			Size   int64  `json:"size"`
			Href   string `json:"href"`
		} `json:"files"`
		Authors []string `json:"authors"`
	} `json:"book"`
}

func (b Book) String() string {
	var size int64
	for _, file := range b.Book.FILES {
		size = size + file.Size
	}
	return fmt.Sprintf("%s-- [%s] %s", b.Book.Title, strings.Join(b.Book.Authors, ","), humanize.Bytes(uint64(size)))
}
func (tale *TaleBook) Next() (*Book, error) {
	tale.index++
	var api = urlJoin(tale.api, "api", "book", fmt.Sprintf("%d", tale.index))
	if err := tale.check(api); err != nil {
		return nil, err
	}
	response, err := tale.client.Get(api)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var book Book
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&book); err != nil {
		return nil, err
	}
	if book.ERROR != "ok" {
		return nil, fmt.Errorf("%s %s", api, book.ERROR)
	}
	return &book, nil
}

func (tale *TaleBook) Download(b *Book, dir string) error {
	for _, file := range b.Book.FILES {
		response, err := tale.client.Get(urlJoin(tale.api, file.Href))
		if err != nil {
			return err
		}
		defer response.Body.Close()
		file := filepath.Join(dir, filename(response))
		fh, err := os.Create(file)
		if err != nil {
			return err
		}
		_, err = io.Copy(fh, response.Body)
		if err != nil {
			fh.Close()
			return err
		}
		fh.Close()
	}
	return nil
}

func (tale *TaleBook) check(api string) error {
	response, err := tale.client.Get(api)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(" %s", response.Status)
	}
	return nil
}

func NewTableBook(site string, opstions ...func(*TaleBook)) (*TaleBook, error) {
	var client http.Client = http.Client{
		Timeout: time.Duration(30) * time.Second,
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client.Jar = jar
	tb := &TaleBook{
		api:    site,
		client: &client,
	}
	for _, option := range opstions {
		option(tb)
	}
	return tb, nil
}

func WithTimeOutOption(timeout time.Duration) func(*TaleBook) {
	return func(tb *TaleBook) {
		tb.client.Timeout = timeout
	}
}

func WithLogin(user string, password string) func(*TaleBook) {
	return func(tb *TaleBook) {
		api := urlJoin(tb.api, "api/user/sign_in")
		respnose, err := tb.client.PostForm(api, map[string][]string{
			"username": []string{user},
			"password": []string{password},
		})
		tb.err = err
		defer respnose.Body.Close()
	}
}
