package cmd

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora/v4"
	"github.com/spf13/cobra"

	"github.com/durandj/git-ignore/internal"
)

func newGenerateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generates a .gitignore file",
		Long:  "Generates a .gitignore file based on certain applications or options",
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

			contents, err := client.Generate(args)

			if err != nil {
				fmt.Println(
					aurora.Sprintf(
						aurora.Red("Unable generate gitignore file\n%s"),
						err,
					),
				)
				os.Exit(1)
			}

			fmt.Println(contents)
		},
	}
}
