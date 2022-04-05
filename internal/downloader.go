package internal

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// downloadWorker is the download instance.
type downloadWorker struct {
	website      string
	wait         *sync.WaitGroup
	progress     *storage
	client       *http.Client
	userAgent    string
	retry        int
	downloadPath string
	formats      []string
	rename       bool
}

// Download would start download books from given website.
func (worker *downloadWorker) Download() {
	// Try to acquire book ID from storage.
	for bookID := worker.progress.AcquireBookID(); bookID != NoBookToDownload; bookID = worker.progress.AcquireBookID() {
		// Acquire book info.
		var info *BookResponse
		for i := 0; i < worker.retry; i++ {
			var err error

			info, err = worker.queryBookInfo(bookID)
			if err == nil {
				break
			}

			// Log the error after last try.
			if i == worker.retry-1 {
				log.Fatalln(err)
			}
		}

		if info == nil {
			log.Printf("Book with ID [%d] is not exist on target website", bookID)
			worker.downloadedBook(bookID)
			continue
		}

		// Find formats to download.
		for _, file := range info.Book.Files {
			for i := 0; i < worker.retry; i++ {
				var err error

				err = worker.downloadBook(bookID, info.Book.Title, file.Format, file.Href)
				if err == nil {
					break
				}

				// Log the error after last try.
				if i == worker.retry-1 {
					log.Fatalln(err)
				}
			}
		}

		worker.downloadedBook(bookID)
	}

	// Finish this download worker
	worker.wait.Done()
}

// downloadedBook would record the download statue into storage.
func (worker *downloadWorker) downloadedBook(bookID int64) {
	if err := worker.progress.SaveBookID(bookID); err != nil {
		log.Fatalln(err)
	}
}

// downloadBook will download the book file from
func (worker *downloadWorker) downloadBook(bookID int64, title, format, href string) error {
	valid := false
	for _, f := range worker.formats {
		if f == format {
			valid = true
			break
		}
	}

	if !valid {
		// Skip this format.
		return nil
	}

	// Start download.
	site := GenerateUrl(worker.website, href)
	req, err := http.NewRequest(http.MethodGet, site, http.NoBody)
	if err != nil {
		return fmt.Errorf("illegal book download request: %w", err)
	}

	req.Header.Set("User-Agent", worker.userAgent)
	resp, err := worker.client.Do(req)
	if err != nil {
		return WrapTimeOut(err)
	}

	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	// Generate file name.
	filename := strconv.FormatInt(bookID, 10) + "." + strings.ToLower(format)
	// Use readable name.
	if !worker.rename {
		name := Filename(resp)
		if name == "" {
			filename = title + "." + strings.ToLower(format)
		} else {
			filename = name
		}
	}

	// Generate the file path.
	file := filepath.Join(worker.downloadPath, filename)

	// Remove the exist file.
	if _, err := os.Stat(file); err == nil {
		if err := os.Remove(file); err != nil {
			return err
		}
	}

	// Create file writer.
	writer, err := os.Create(file)
	if err != nil {
		return err
	}
	defer func() { _ = writer.Close() }()

	// Write file content
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return WrapTimeOut(err)
	}

	return nil
}

// queryBookInfo will find the required book information.
func (worker *downloadWorker) queryBookInfo(bookID int64) (*BookResponse, error) {
	site := GenerateUrl(worker.website, "/api/book", strconv.FormatInt(bookID, 10))

	req, err := http.NewRequest(http.MethodGet, site, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", worker.userAgent)

	resp, err := worker.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	result := &BookResponse{}
	if err := DecodeResponse(resp, result); err != nil {
		return nil, err
	}

	switch result.Err {
	case SuccessStatus:
		return result, nil
	case BookNotFound:
		return nil, nil
	default:
		return nil, errors.New(result.Msg)
	}
}
