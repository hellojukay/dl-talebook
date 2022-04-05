package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

// Used for download books from talebook website.
var (
	downloadPath, _ = os.Getwd()       // Use the executed directory as the default download path.
	cookieFile      = "cookies"        // The cookie file to use in this download progress.
	progressFile    = "progress"       // The progress file serving the remaining book id.
	initialBookID   = 1                // The book id start to download.
	formats         []string           // The file formats you want to download
	threads         = 1                // The concurrent goroutine counts.
	timeout         = 10 * time.Second // The request timeout for a single request.
	retry           = 5                // The maximum retry times for a timeout request.
	// The user agent for the download request.
	userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.4844.51 Safari/537.36"
)

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
		if threads < 1 {
			log.Fatalln(ErrThreadNumber)
		}
		if initialBookID < 1 {
			log.Fatalln(ErrInitialBookID)
		}
		if retry < 1 {
			log.Fatalln(ErrRetryTimes)
		}

		startDownload()
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Add flags for use info.
	downloadCmd.Flags().StringVarP(&website, "website", "w", "", "The talebook website.")
	downloadCmd.Flags().StringVarP(&username, "username", "u", "", "The account login name.")
	downloadCmd.Flags().StringVarP(&password, "password", "p", "", "The account password.")
	downloadCmd.Flags().StringVarP(&downloadPath, "download", "d", downloadPath, "The book directory you want to use, default would be current working directory.")
	downloadCmd.Flags().StringVarP(&cookieFile, "cookie", "c", cookieFile, "The cookie file name you want to use, it would be saved under the download directory.")
	downloadCmd.Flags().StringVarP(&progressFile, "progress", "g", progressFile, "The download progress file name you want to use, it would be saved under the download directory.")
	downloadCmd.Flags().IntVarP(&initialBookID, "initial", "i", initialBookID, "The book id you want to start download. It should exceed 0.")
	downloadCmd.Flags().StringSliceVarP(&formats, "format", "f", []string{"EPUB", "MOBI", "PDF"}, "The file formats you want to download.")
	downloadCmd.Flags().IntVarP(&threads, "thread", "t", threads, "The number of concurrent download request.")
	downloadCmd.Flags().DurationVarP(&timeout, "timeout", "o", timeout, "The max pending time for download request.")
	downloadCmd.Flags().IntVarP(&retry, "retry", "r", retry, "The max retry times for timeout download request.")
	downloadCmd.Flags().StringVarP(&userAgent, "user-agent", "a", userAgent, "Set User-Agent for download request.")

	_ = downloadCmd.MarkFlagRequired("website")
}

// startDownload would start download books from the given configuration.
func startDownload() {
}
