/*
2024 David Rose-Franklin <david.franklin.dev@gmail.ocm>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/DavidRR-F/jenklog/internal/config"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth [arg]",
	Short: "Initialize authentication credentials",
	Args:  cobra.ExactArgs(1),
	Run:   executeAuth,
}

func executeAuth(cmd *cobra.Command, args []string) {
	url := args[0]

	user, err := cmd.Flags().GetString("user")
	if err != nil {
		fmt.Printf("Error getting user flag: %v\n", err)
		os.Exit(1)
	}

	token, err := cmd.Flags().GetString("token")
	if err != nil {
		fmt.Printf("Error getting token flag: %v\n", err)
		os.Exit(1)
	}

	if err := config.SaveConfig(user, token, url); err != nil {
		fmt.Printf("Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nConfiguration saved successfully in ~/.jenklog-config")
}

func init() {
	authCmd.Flags().StringP("user", "u", "", "Jenkins Username")
	authCmd.Flags().StringP("token", "t", "", "Jenkins API Token")

	authCmd.MarkFlagRequired("username")
	authCmd.MarkFlagRequired("token")
	rootCmd.AddCommand(authCmd)
}
