package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
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
	if book.Err != "ok" {
		return nil, fmt.Errorf("%s %s", api, book.Err)
	}
	return &book, nil
}

func (tale *TaleBook) Download(b *Book, dir string) error {
	for _, file := range b.Book.Files {
		downloadURL := urlJoin(tale.api, file.Href)
		response, err := tale.client.Get(downloadURL)
		if err != nil {
			return wrapperTimeOutError(err)
		}
		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("%s %s", downloadURL, response.Status)
		}
		defer response.Body.Close()
		name := filename(response)
		if name == "" {
			name = b.Book.Title + "." + strings.ToLower(file.Format)
		}
		filepath := filepath.Join(dir, name)
		if info, err := os.Stat(filepath); err == nil {
			if file.Size == info.Size() {
				return os.ErrExist
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
			return nil, err
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
func WithLoginOption(user string, password string) func(*TaleBook) {
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

func (tb *TaleBook) getInfo() {

	api := urlJoin(tb.api, "api/user/info")
	req, err := http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		tb.err = wrapperTimeOutError(err)
		return
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
	var info ServerInfo
	decoder := json.NewDecoder(respnose.Body)
	if err = decoder.Decode(&info); err != nil {
		tb.err = err
		return
	}
	tb.ServerInfo = info
}
