package cmd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	goconfluence "github.com/virtomize/confluence-go-api"
	"golang.org/x/sync/errgroup"
)

var debugCmd = &cobra.Command{
	Use:   "debug <parent_page_id> <count>",
	Short: "",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {
		parentID := args[0]
		count, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}

		client, err := getClient()
		if err != nil {
			return err
		}

		parent, err := client.GetContentByID(context.Background(), parentID)
		if err != nil {
			return err
		}

		fmt.Printf("count : %d\n", count)
		fmt.Printf("parent: [%s] %s\n", parent.ID, parent.Title)
		fmt.Print("Are you sure ? (yes/no): ")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		if input.Text() != "yes" {
			return errors.New("aborted")
		}

		eg, ctx := errgroup.WithContext(context.Background())

		for i := 0; i < count; i++ {
			eg.Go(func() error {
				// FIXME: CreatePage() doees not support context cancel.
				return client.CreatePage(ctx, &goconfluence.Content{
					Type:  "page",
					Title: uuid.New().String(),
					Ancestors: []goconfluence.Ancestor{
						{
							ID: parentID,
						},
					},
					Body: goconfluence.Body{
						Storage: goconfluence.Storage{
							Value:          "#api-test\nnew sub\npage",
							Representation: "storage",
						},
					},
					Version: &goconfluence.Version{
						Number: 1,
					},
					Space: &goconfluence.Space{
						Key: parent.Space.Key,
					},
				})
			})
		}

		if err := eg.Wait(); err != nil {
			return err
		}

		return nil
	},
	Args: cobra.ExactArgs(2),
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
