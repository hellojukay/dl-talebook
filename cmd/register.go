package cmd

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hellojukay/dl-talebook/internal"
)

// Used for register account on talebook website.
type registerConfig struct {
	website   string
	username  string
	password  string
	email     string
	userAgent string
}

// Arguments instance.
var regConf = registerConfig{}

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register account on talebook.",
	Long: `Some talebook website need a user account for downloading books.
You can use this register command for creating account.`,
	Run: func(cmd *cobra.Command, args []string) {
		register()
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	// Add flags for use info.
	registerCmd.Flags().StringVarP(&regConf.website, "website", "w", "", "The talebook website.")
	registerCmd.Flags().StringVarP(&regConf.username, "username", "u", "", "The account login name.")
	registerCmd.Flags().StringVarP(&regConf.password, "password", "p", "", "The account password.")
	registerCmd.Flags().StringVarP(&regConf.email, "email", "e", "", "The account email.")
	registerCmd.Flags().StringVarP(&regConf.userAgent, "user-agent", "a", internal.DefaultUserAgent, "The account email.")

	_ = registerCmd.MarkFlagRequired("website")
	_ = registerCmd.MarkFlagRequired("username")
	_ = registerCmd.MarkFlagRequired("password")
	_ = registerCmd.MarkFlagRequired("email")
}

// register will create account on given website
func register() {
	website := internal.GenerateUrl(regConf.website, "/api/user/sign_up")
	referer := internal.GenerateUrl(regConf.website, "/signup")
	values := url.Values{
		"username": {regConf.username},
		"password": {regConf.password},
		"nickname": {regConf.username},
		"email":    {regConf.email},
	}

	req, err := http.NewRequest(http.MethodPost, website, strings.NewReader(values.Encode()))
	if err != nil {
		log.Fatalln("Illegal login request: %w", err)
	}
	req.Header.Set("User-Agent", regConf.userAgent)
	req.Header.Set("referer", referer)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	form, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() { _ = form.Body.Close() }()
	if form.StatusCode != http.StatusOK {
		log.Fatalf("Error in register user, message: %s", form.Status)
	}

	result := &internal.CommonResponse{}
	if err = internal.DecodeResponse(form, result); err != nil {
		log.Fatalln(err)
	}

	if result.Err == internal.SuccessStatus {
		log.Printf("Register success.")
	} else {
		log.Fatalf("Register failed, reason: %s", result.Err)
	}
}
