package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ymtdzzz/confl-mv/client"
)

var (
	domain   string
	username string
	apiKey   string
)

var rootCmd = &cobra.Command{
	Use:   "confl-mv",
	Short: "",
	Long:  "",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "user_name@example.com", "Your Confluence username")
	rootCmd.PersistentFlags().StringVarP(&apiKey, "apikey", "a", "your api key", "Your Confluence api key")
	rootCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "your_site.atlassian.net", "Your Confluence domain")
	rootCmd.MarkFlagRequired("username")
	rootCmd.MarkFlagRequired("apikey")
	rootCmd.MarkFlagRequired("domain")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getClient() (client.Client, error) {
	return client.NewClient(
		fmt.Sprintf("https://%s/wiki/rest/api", domain),
		fmt.Sprintf("https://%s/wiki/api/v2", domain),
		username,
		apiKey,
	)
}
