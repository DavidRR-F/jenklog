/*
2024 David Rose-Franklin <david.franklin.dev@gmail.ocm>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jenklog",
	Short: "Query Jenkins Logs",
	Long: `
  This tool is a more verbose implementation of the jenkins-cli's console subcommand that allows 
  users to query and parse multiple Jenkins Logs simultaneously throught the terminal:

  Example:

  jenklog job <job-name> --build last --count 0 --stage Deploy --filter-status success
  `,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
