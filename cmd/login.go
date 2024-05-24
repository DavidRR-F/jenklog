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

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Jenkins login",
	Long:  "Jenkins login command to save credentials to a config file",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		token, _ := cmd.Flags().GetString("token")
		url, _ := cmd.Flags().GetString("url")

		if err := config.SaveConfig(username, token, url); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Configuration saved successfully")
	},
}

func init() {
	loginCmd.Flags().StringP("username", "u", "", "Jenkins username")
	loginCmd.Flags().StringP("token", "t", "", "Jenkins API token")
	loginCmd.Flags().StringP("url", "r", "", "Jenkins URL")

	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("token")
	loginCmd.MarkFlagRequired("url")

	rootCmd.AddCommand(loginCmd)
}
