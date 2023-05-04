package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Execute runs the root command.
func Execute() {
	rootCmd := &cobra.Command{
		Use:   "git-ignore",
		Short: ".gitignore generator",
		Long:  "Generates contents for a .gitignore file using gitignore.io",
	}

	rootCmd.AddCommand(
		newGenerateCommand(),
		newListCommand(),
		newUpdateCommand(),
		newVersionCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
