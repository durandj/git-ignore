package cmd

import (
	"fmt"

	"github.com/durandj/git-ignore/pkg/gitignore"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current version of git-ignore",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(gitignore.VERSION)
	},
}
