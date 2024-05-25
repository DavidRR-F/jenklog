/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
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
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		username, _ := cmd.Flags().GetString("username")
		token, _ := cmd.Flags().GetString("token")

		if err := config.SaveConfig(username, token, url); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\nConfiguration saved successfully in ~/.jenklog-config")
	},
}

func init() {
	authCmd.Flags().StringP("user", "u", "", "Jenkins Username")
	authCmd.Flags().StringP("token", "t", "", "Jenkins API Token")

	authCmd.MarkFlagRequired("username")
	authCmd.MarkFlagRequired("token")
	rootCmd.AddCommand(authCmd)
}
