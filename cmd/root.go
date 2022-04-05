package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var (
	website  = ""
	username = ""
	password = ""
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dl-talebook",
	Short: "A command line base downloader for downloading books from talebook server.",
	Long: `You can use this command to register account and download book.
The url for talebook should be provided, the formats is also
optional.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
