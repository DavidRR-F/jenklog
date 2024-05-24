/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/DavidRR-F/jenklog/internal/jenkins"
)

// jobCmd represents the job command
var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "A brief description of your command",
	Long:  `A longer description`,
	Args:  cobra.ExactArgs(1),
	Run:   getJobLogs,
}

func getJobLogs(cmd *cobra.Command, args []string) {
	job := args[0]
	stage, _ := cmd.Flags().GetString("stage")
	builds, _ := cmd.Flags().GetString("builds")

	statusMap := map[string]string{
		"success":  "SUCCESS",
		"failure":  "FAILURE",
		"aborted":  "ABORTED",
		"unstable": "UNSTABLE",
		"notbuilt": "NOT_BUILT",
	}

	//TODO: Parse builds
	// For each build get run logs

	err := jenkins.GetJenkinsJobLogs(job, stage)

	if err != nil {
		fmt.Println(err)
	}
}

func parseStage() error {
	return nil
}

func init() {
	jobCmd.Flags().StringP("stage", "s", "", "Job Build Stage Name")
	jobCmd.Flags().StringP("builds", "b", "last", "Job Build Number")
	jobCmd.Flags().StringP("status", "bs", "", "Job Build Status")

	rootCmd.AddCommand(jobCmd)
}
