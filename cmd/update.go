package cmd

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"

	"github.com/durandj/git-ignore/internal"
)

func newUpdateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Updates stored data",
		Long:  "Ensures that any stored data is updated",
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

			err = client.Update()
			if err != nil {
				fmt.Println(
					aurora.Sprintf(
						aurora.Red("Error while updating\n%s"),
						err,
					),
				)
				os.Exit(1)
			}

			fmt.Println(aurora.Green("Update complete!"))
		},
	}
}
