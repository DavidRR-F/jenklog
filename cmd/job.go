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

	//TODO: Parse builds
	// For each build get run logs
	// Two Use Cases
	// 1. Grabbing Single build log skip fetching runs
	// 2. Grabbing Multiple build logs parse builds param

	runs, err := jenkins.GetJenkinsJobRuns(job)
	if err != nil {
		fmt.Println(err)
	}
	var logs []jenkins.Log
	for _, run := range runs {
		log, err := getJobLog(job, stage, run.ID)

		if err != nil {
			fmt.Println("Cant Find Log")
		}

		logs = append(logs, log)
	}

	// print logs
}

func getJobRuns(job string) {

}

func getJobLog(job, stage, build string) (jenkins.Log, error) {
	log, err := jenkins.GetJenkinsJobLog(job, build)

	if err != nil {
		fmt.Println(err)
	}

	if stage != "" {
		err := log.ParseByStage(stage)
		if err != nil {
			fmt.Println("Stage not found")

		}
	}
}

func init() {
	jobCmd.Flags().StringP("stage", "s", "", "Job Build Stage Name")
	jobCmd.Flags().StringP("builds", "b", "last", "Job Build Number")
	jobCmd.Flags().StringP("status", "o", "", "Job Build Status")

	rootCmd.AddCommand(jobCmd)
}
