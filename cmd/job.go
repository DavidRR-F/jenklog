/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"

	"github.com/DavidRR-F/jenklog/internal/jenkins"
)

const (
	SUCCESS   = "SUCCESS"
	FAILURE   = "FAILURE"
	ABORTED   = "ABORTED"
	UNSTABLE  = "UNSTABLE"
	NOT_BUILT = "NOT_BUILT"
)

var statusMap = map[string]string{
	"success":  SUCCESS,
	"failure":  FAILURE,
	"aborted":  ABORTED,
	"unstable": UNSTABLE,
	"notbuilt": NOT_BUILT,
}

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
	build, _ := cmd.Flags().GetString("build")
	count, _ := cmd.Flags().GetInt("count")
	filter, _ := cmd.Flags().GetString("status-filter")

	//valifate filter

	if count < 1 {
		if filter != "" {
			fmt.Println("count must be greater than 0 to filer by status")
			os.Exit(1)
		}
		log, err := getJobLog(job, stage, build)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		log.Print()
	} else {
		//multi logs
		runs, err := getJobRuns(job, build, filter, count)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		logChannel := make(chan jenkins.Log, len(runs))
		errChannel := make(chan error, len(runs))
		var wg sync.WaitGroup
		for _, run := range runs {
			wg.Add(1)
			go func(job, stage string, run jenkins.Run) {
				defer wg.Done()
				log, err := getJobLog(job, stage, run.ID)
				if err != nil {
					errChannel <- err
				}
				logChannel <- log
			}(job, stage, run)
		}
		go func() {
			wg.Wait()
			close(logChannel)
			close(errChannel)
		}()

		for log := range logChannel {
			log.Print()
		}
		for err := range errChannel {
			fmt.Println(err)
		}
	}
}

func getJobRuns(job, build, filter string, count int) ([]jenkins.Run, error) {
	runs, err := jenkins.GetJenkinsJobRuns(job)

	if err != nil {
		return []jenkins.Run{}, err
	}

	if build == "last" {
		runs = runs[:count+1]
	} else {
		runs, err = sliceRuns(runs, build, count)
		if err != nil {
			return []jenkins.Run{}, err
		}
	}
	if filter != "" {
		filter, valid := statusMap[filter]
		if !valid {
			return []jenkins.Run{}, fmt.Errorf("invalid status filter: %s", filter)
		}
		var filteredRuns []jenkins.Run
		for _, run := range runs {
			if run.Status == filter {
				filteredRuns = append(filteredRuns, run)
			}
		}
		return filteredRuns, nil
	}
	return runs, nil
}

func sliceRuns(runs []jenkins.Run, build string, count int) ([]jenkins.Run, error) {
	startIndex := -1
	for i, run := range runs {
		if run.ID == build {
			startIndex = i
			break
		}
	}

	if startIndex == -1 {
		return []jenkins.Run{}, fmt.Errorf("value %s not found in the builds", build)
	}

	endIndex := startIndex + count + 1
	if endIndex > len(runs) {
		endIndex = len(runs)
	}

	return runs[startIndex:endIndex], nil
}

func getJobLog(job, stage, build string) (jenkins.Log, error) {
	log, err := jenkins.GetJenkinsJobLog(job, build)

	if err != nil {
		return jenkins.Log{}, err
	}

	if stage != "" {
		err := log.ParseByStage(stage)
		if err != nil {
			return jenkins.Log{}, err
		}
	}
	return log, nil
}

func init() {
	jobCmd.Flags().StringP("build", "b", "last", "Job Build Number")
	jobCmd.Flags().StringP("status-filter", "f", "", "Filter Job Builds By Status")
	jobCmd.Flags().StringP("stage", "s", "", "Job Build Stage Name")
	jobCmd.Flags().IntP("count", "c", 0, "Job Build Log")

	rootCmd.AddCommand(jobCmd)
}
