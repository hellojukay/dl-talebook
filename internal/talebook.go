package internal

import (
	"log"
	"os"
	"time"
)

// The download config.
type config struct {
	Website       string        // The website for talebook.
	Username      string        // The login user.
	Password      string        // The password for login user.
	DownloadPath  string        // Use the executed directory as the default download path.
	CookieFile    string        // The cookie file to use in this download progress.
	ProgressFile  string        // The progress file serving the remaining book id.
	InitialBookID int           // The book id start to download.
	Formats       []string      // The file formats you want to download
	Threads       int           // The concurrent goroutine counts.
	Timeout       time.Duration // The request timeout for a single request.
	Retry         int           // The maximum retry times for a timeout request.
	UserAgent     string        // The user agent for the download request.
}

// The main instance for start downloading the book.
type talebook struct {
	website    string
	progress   *storage
	login      *loginWorker
	downloader *downloadWorker
	channel    chan int64
}

// The login instance.
type loginWorker struct {
	username   string
	password   string
	cookieFile string
	userAgent  string
}

// The download instance.
type downloadWorker struct {
	formats      []string
	timeout      time.Duration
	retry        int
	downloadPath string
	cookieFile   string
	userAgent    string
}

// NewDownloadConfig will return a default blank config.
func NewDownloadConfig() *config {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	return &config{
		DownloadPath:  dir,
		CookieFile:    "cookies",
		ProgressFile:  "progress",
		InitialBookID: 1,
		Formats:       []string{"EPUB", "MOBI", "PDF"},
		Threads:       1,
		Timeout:       10 * time.Second,
		Retry:         5,
		UserAgent:     "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36",
	}
}

// NewTalebook will create the download instance.
func NewTalebook(c *config) *talebook {
	return &talebook{}
}

func (*talebook) Start() {

}
