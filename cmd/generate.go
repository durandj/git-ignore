package cmd

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"

	"github.com/durandj/git-ignore/pkg/gitignore"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a .gitignore file",
	Long:  "Generates a .gitignore file based on certain applications or options",
	Run: func(cmd *cobra.Command, args []string) {
		client := gitignore.Client{}
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
