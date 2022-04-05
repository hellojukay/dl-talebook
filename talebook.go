package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

var (
	NO_MORE_BOOK_ERROR = errors.New("there is no more books")
)

type ServerInfo struct {
	Err string `json:"err"`
	Cdn string `json:"cdn"`
	Sys struct {
		Books      int    `json:"books"`
		Tags       int    `json:"tags"`
		Authors    int    `json:"authors"`
		Publishers int    `json:"publishers"`
		Series     int    `json:"series"`
		Mtime      string `json:"mtime"`
		Users      int    `json:"users"`
		Active     int    `json:"active"`
		Version    string `json:"version"`
		Title      string `json:"title"`
		Socials    []struct {
			Text  string `json:"text"`
			Value string `json:"value"`
			Help  bool   `json:"help"`
			Link  string `json:"link"`
		} `json:"socials"`
		Friends []struct {
			Text string `json:"text"`
			Href string `json:"href"`
		} `json:"friends"`
		Footer string `json:"footer"`
		Allow  struct {
			Register bool `json:"register"`
			Download bool `json:"download"`
			Push     bool `json:"push"`
			Read     bool `json:"read"`
		} `json:"allow"`
	} `json:"sys"`
	User struct {
		Avatar      string `json:"avatar"`
		IsLogin     bool   `json:"is_login"`
		IsAdmin     bool   `json:"is_admin"`
		Nickname    string `json:"nickname"`
		Email       string `json:"email"`
		KindleEmail string `json:"kindle_email"`
		Extra       struct {
		} `json:"extra"`
	} `json:"user"`
	Msg string `json:"msg"`
}
type TaleBook struct {
	api        string
	index      int
	client     *http.Client
	err        error
	userAgent  string
	cookie     string
	verbose    bool
	ServerInfo ServerInfo
}

type Book struct {
	Err          string `json:"err"`
	KindleSender string `json:"kindle_sender"`
	Book         struct {
		ID            int         `json:"id"`
		Title         string      `json:"title"`
		Rating        int         `json:"rating"`
		CountVisit    int         `json:"count_visit"`
		CountDownload int         `json:"count_download"`
		Timestamp     string      `json:"timestamp"`
		Pubdate       string      `json:"pubdate"`
		Collector     string      `json:"collector"`
		Authors       []string    `json:"authors"`
		Author        string      `json:"author"`
		Tags          []string    `json:"tags"`
		AuthorSort    string      `json:"author_sort"`
		Publisher     string      `json:"publisher"`
		Comments      string      `json:"comments"`
		Series        interface{} `json:"series"`
		Language      interface{} `json:"language"`
		Isbn          string      `json:"isbn"`
		Files         []struct {
			Format string `json:"format"`
			Size   int64  `json:"size"`
			Href   string `json:"href"`
		} `json:"files"`
		IsPublic bool   `json:"is_public"`
		IsOwner  bool   `json:"is_owner"`
		Img      string `json:"img"`
	} `json:"book"`
	Msg string `json:"msg"`
}

func (b Book) String() string {
	var size int64
	for _, file := range b.Book.Files {
		size = size + file.Size
	}
	return fmt.Sprintf("%s-- [%s] %s", b.Book.Title, strings.Join(b.Book.Authors, ","), humanize.Bytes(uint64(size)))
}
func (tale *TaleBook) Next() (*Book, error) {
	tale.index++
	if tale.index > tale.ServerInfo.Sys.Books {
		return nil, NO_MORE_BOOK_ERROR
	}
	var api = urlJoin(tale.api, "api", "book", fmt.Sprintf("%d", tale.index))
	if tale.verbose {
		log.Printf("feth book from %s", api)
	}
	req, err := http.NewRequest(http.MethodGet, api, nil)
	if tale.userAgent != "" {
		req.Header.Set("user-agent", tale.userAgent)
	}
	if tale.cookie != "" {
		req.Header.Set("cookie", tale.cookie)
	}
	response, err := tale.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", response.Status)
	}
	var book Book
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&book); err != nil {
		return nil, fmt.Errorf("parse json failed %w", err)
	}
	if book.Err != "ok" {
		return nil, fmt.Errorf("%s %s", api, book.Err)
	}
	return &book, nil
}

func (tale *TaleBook) Download(b *Book, dir string) error {
	for _, file := range b.Book.Files {
		downloadURL := urlJoin(tale.api, file.Href)
		req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
		if err != nil {
			return err
		}
		if tale.cookie != "" {
			req.Header.Set("cookie", tale.cookie)
		}
		if tale.userAgent != "" {
			req.Header.Set("user-agent", tale.userAgent)
		}
		response, err := tale.client.Do(req)
		if err != nil {
			return wrapperTimeOutError(err)
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("%s %s", downloadURL, response.Status)
		}

		name := filename(response)
		if name == "" {
			name = b.Book.Title + "." + strings.ToLower(file.Format)
		}
		filepath := filepath.Join(dir, name)
		if info, err := os.Stat(filepath); err == nil {
			if file.Size == info.Size() {
				return fmt.Errorf("%s %w", filepath, os.ErrExist)
			} else {
				log.Printf("expected file size %d, actual file size %d, so removing %s, ", file.Size, info.Size(), filepath)
				if err = os.Remove(filepath); err != nil {
					return err
				}
			}
		}
		fh, err := os.Create(filepath)
		if err != nil {
			return err
		}
		_, err = io.Copy(fh, response.Body)
		if err != nil {
			fh.Close()
			return wrapperTimeOutError(err)
		}
		fh.Close()
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
		if tb.err != nil {
			return nil, tb.err
		}
	}

	tb.getInfo()

	return tb, tb.err
}

func WithTimeOutOption(timeout time.Duration) func(*TaleBook) {
	return func(tb *TaleBook) {
		tb.client.Timeout = timeout
	}
}

func WithUserAgentOption(uagent string) func(*TaleBook) {
	return func(tb *TaleBook) {
		tb.userAgent = userAgent
	}
}
func WithUserCookieOption(cookie string) func(*TaleBook) {
	return func(tb *TaleBook) {
		if cookie != "" {
			tb.cookie = cookie
		}
	}
}
func WithLoginOption(user string, password string) func(*TaleBook) {
	return func(tb *TaleBook) {
		if (user != "") && (password != "") {

			data := url.Values{
				"username": []string{user},
				"password": []string{password},
			}
			api := urlJoin(tb.api, "api/user/sign_in")
			req, err := http.NewRequest(http.MethodPost, api, strings.NewReader(data.Encode()))
			if err != nil {
				tb.err = fmt.Errorf("login failed %w", err)
				return
			}

			if tb.userAgent != "" {
				req.Header.Set("user-agent", tb.userAgent)
			}
			req.Header.Set("referer", urlJoin(tb.api, "login"))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			if tb.verbose {
				log.Printf("login %s username: %s password: %s", api, username, password)
			}

			respnose, err := tb.client.Do(req)
			if err != nil {
				tb.err = fmt.Errorf("login failed %w", err)
				return
			}
			defer respnose.Body.Close()
			if respnose.StatusCode != http.StatusOK {
				if tb.verbose {
					io.Copy(os.Stderr, respnose.Body)
				}
				tb.err = fmt.Errorf("%s %s", api, respnose.Status)
				return
			}
			type Result struct {
				Err string `json:"err"`
				Msg string `json:"msg"`
			}
			if tb.verbose {
				log.Printf("%s %s", api, respnose.Status)
			}
			var result Result
			decoder := json.NewDecoder(respnose.Body)
			if err = decoder.Decode(&result); err != nil {
				tb.err = fmt.Errorf("login failed %w", err)
				return
			}
			if result.Err != "ok" {
				tb.err = fmt.Errorf("login failed error: %s,message: %s", result.Err, result.Msg)
				return
			}
			if tb.verbose {
				log.Printf("login %s %s", result.Err, result.Msg)
				return
			}
		}
	}
}

func WithVerboseOption(verbose bool) func(*TaleBook) {
	return func(tb *TaleBook) {
		tb.verbose = verbose
	}
}
func (tb *TaleBook) getInfo() {

	api := urlJoin(tb.api, "api/user/info")
	req, err := http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		tb.err = wrapperTimeOutError(err)
		return
	}
	if tb.cookie != "" {
		req.Header.Set("cookie", tb.cookie)
	}
	if tb.userAgent != "" {
		req.Header.Set("user-agent", tb.userAgent)
	}
	respnose, err := tb.client.Do(req)
	if err != nil {
		tb.err = err
		return
	}
	defer respnose.Body.Close()
	if respnose.StatusCode != http.StatusOK {
		tb.err = fmt.Errorf("%s %s", api, respnose.Status)
		return
	}
	var info ServerInfo
	decoder := json.NewDecoder(respnose.Body)
	if err = decoder.Decode(&info); err != nil {
		tb.err = err
		return
	}
	tb.ServerInfo = info
}
