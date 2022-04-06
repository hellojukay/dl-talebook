package cmd

import (
	"errors"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"

	"github.com/hellojukay/dl-talebook/internal"
)

// Used for downloading books from talebook website.
var downloadConfig = *internal.NewDownloadConfig()

var (
	ErrThreadNumber  = errors.New("illegal thread number, it should exceed 0")
	ErrInitialBookID = errors.New("illegal book id, it should exceed 0")
	ErrRetryTimes    = errors.New("illegal retry times, it should exceed 0")
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download the book from talebook.",
	Run: func(cmd *cobra.Command, args []string) {
		// Validate the download arguments.
		if downloadConfig.Threads < 1 {
			log.Fatalln(ErrThreadNumber)
		}
		if downloadConfig.InitialBookID < 1 {
			log.Fatalln(ErrInitialBookID)
		}
		if downloadConfig.Retry < 1 {
			log.Fatalln(ErrRetryTimes)
		}
		for i, format := range downloadConfig.Formats {
			// Make sure all the format should be upper case.
			downloadConfig.Formats[i] = strings.ToUpper(format)
		}

		// Print download configuration.
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Config Key", "Config Value"})
		cv := reflect.ValueOf(downloadConfig)
		for i := 0; i < cv.NumField(); i++ {
			name := cv.Type().Field(i).Name
			value := cv.Field(i).Interface()
			t.AppendRow([]interface{}{name, value})
		}
		t.Render()

		// Create the downloader
		talebook := internal.NewTalebook(&downloadConfig)

		// Start download books.
		talebook.Start()

		// Finished all the tasks.
		log.Println("Successfully download all the books.")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Add flags for use info.
	downloadCmd.Flags().StringVarP(&downloadConfig.Website, "website", "w", "", "The talebook website.")
	downloadCmd.Flags().StringVarP(&downloadConfig.Username, "username", "u", "", "The account login name.")
	downloadCmd.Flags().StringVarP(&downloadConfig.Password, "password", "p", "", "The account password.")
	downloadCmd.Flags().StringVarP(&downloadConfig.DownloadPath, "download", "d", downloadConfig.DownloadPath,
		"The book directory you want to use, default would be current working directory.")
	downloadCmd.Flags().StringVarP(&downloadConfig.CookieFile, "cookie", "c", downloadConfig.CookieFile,
		"The cookie file name you want to use, it would be saved under the download directory.")
	downloadCmd.Flags().StringVarP(&downloadConfig.ProgressFile, "progress", "g", downloadConfig.ProgressFile,
		"The download progress file name you want to use, it would be saved under the download directory.")
	downloadCmd.Flags().IntVarP(&downloadConfig.InitialBookID, "initial", "i", downloadConfig.InitialBookID,
		"The book id you want to start download. It should exceed 0.")
	downloadCmd.Flags().StringSliceVarP(&downloadConfig.Formats, "format", "f", downloadConfig.Formats,
		"The file formats you want to download.")
	downloadCmd.Flags().IntVarP(&downloadConfig.Threads, "thread", "t", downloadConfig.Threads, "The number of concurrent download request.")
	downloadCmd.Flags().DurationVarP(&downloadConfig.Timeout, "timeout", "o", downloadConfig.Timeout,
		"The max pending time for download request.")
	downloadCmd.Flags().IntVarP(&downloadConfig.Retry, "retry", "r", downloadConfig.Retry, "The max retry times for timeout download request.")
	downloadCmd.Flags().StringVarP(&downloadConfig.UserAgent, "user-agent", "a", downloadConfig.UserAgent,
		"Set User-Agent for download request.")
	downloadCmd.Flags().BoolVarP(&downloadConfig.Rename, "rename", "n", downloadConfig.Rename, "Rename the book file by book ID.")

	_ = downloadCmd.MarkFlagRequired("website")
}
