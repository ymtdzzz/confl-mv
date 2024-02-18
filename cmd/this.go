package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	goconfluence "github.com/virtomize/confluence-go-api"
	"golang.org/x/sync/errgroup"
)

var thisCmd = &cobra.Command{
	Use:   "this <target_page_id> <destination_page_id>",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		t, d := args[0], args[1]

		client, err := getClient()
		if err != nil {
			return err
		}

		eg, ctx := errgroup.WithContext(context.Background())
		var (
			target      *goconfluence.Content
			destination *goconfluence.Content
		)
		eg.Go(func() error {
			target, err = client.GetContentByID(ctx, t)
			return err
		})
		eg.Go(func() error {
			destination, err = client.GetContentByID(ctx, d)
			return err
		})

		if err := eg.Wait(); err != nil {
			return err
		}

		fmt.Printf("target     : [%s] %s\n", target.ID, target.Title)
		fmt.Printf("destination: [%s] %s\n", destination.ID, destination.Title)
		fmt.Print("Are you sure ? (yes/no): ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		if input.Text() != "yes" {
			return errors.New("aborted")
		}

		ok, err := client.MovePage(context.Background(), t, d)
		if err != nil {
			return err
		}

		if ok {
			fmt.Printf("Page [%s] moved to [%s]\n", target.ID, destination.ID)
			return nil
		}
		fmt.Printf(`Failed to move page [%s] to [%s]
because it has over 99 child pages.
Use child command to move this page`, target.ID, destination.ID)

		return nil
	},
	Args: cobra.ExactArgs(2),
}

func init() {
	rootCmd.AddCommand(thisCmd)
}
