package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"

	"github.com/durandj/git-ignore/internal"
)

func newListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Gets a list of all possible gitignore options",
		Long:  "Retrieves a list of all the options that can be specified for creating a .gitignore file",
		Run: func(cmd *cobra.Command, args []string) {
			client, err := internal.NewClient()
			if err != nil {
				fmt.Println(
					aurora.Sprintf(
						aurora.Red("Error creating client\n%s"),
						err,
					),
				)
				os.Exit(1)
			}

			options, err := client.List()
			if err != nil {
				fmt.Println(
					aurora.Sprintf(
						aurora.Red("Error retrieving list of options:\n%s"),
						err,
					),
				)
				os.Exit(1)
			}

			fmt.Println(aurora.Bold("Options:"))
			fmt.Println(strings.Join(options, ", "))
		},
	}
}
