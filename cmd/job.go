/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/spf13/cobra"

	"github.com/DavidRR-F/jenklog/internal/jenkins"
)

var jobCmd = &cobra.Command{
	Use:   "job [jobName]",
	Short: "Query jenkins job logs",
	Args:  cobra.ExactArgs(1),
	Run:   getJobLogs,
}

func getJobLogs(cmd *cobra.Command, args []string) {
	job := args[0]
	stage, _ := cmd.Flags().GetString("stage")
	build, _ := cmd.Flags().GetString("build")
	count, _ := cmd.Flags().GetInt("prev-count")
	filter, _ := cmd.Flags().GetString("filter-status")

	if !jenkins.IsValidBuildOption(build) {
		fmt.Println("Not a valid build option")
		os.Exit(1)
	}

	if count < 1 {
		if filter != "" {
			fmt.Println("must select multiple builds to fitler by status")
			os.Exit(1)
		}
		log, err := getJobLog(job, stage, build)
		if err != nil {
			fmt.Println(err)
		}
		log.Print()
	} else {
		runs, err := getJobRuns(job, build, filter, count)
		if err != nil {
			fmt.Println(err)
		}
		logs := make([]jenkins.Log, len(runs))
		var wg sync.WaitGroup
		for i, run := range runs {
			wg.Add(1)
			go func(job, stage string, run jenkins.Run, i int) {
				defer wg.Done()
				log, err := getJobLog(job, stage, run.ID)
				if err != nil {
					fmt.Println(err)
					return
				}
				logs[i] = log
			}(job, stage, run, i)
		}
		wg.Wait()

		for _, log := range logs {
			log.Print()
		}
	}
}

func getJobRuns(job, build, filter string, count int) ([]jenkins.Run, error) {
	if _, err := strconv.Atoi(build); err != nil && build != jenkins.LAST_BUILD {
		buildInfo, err := jenkins.GetJenkinsJobInfo(job, build)
		if err != nil {
			return []jenkins.Run{}, err
		}
		build = buildInfo.ID
	}

	runs, err := jenkins.GetJenkinsJobRuns(job)

	if err != nil {
		return []jenkins.Run{}, err
	}

	if build == jenkins.LAST_BUILD {
		if count+1 < len(runs) {
			runs = runs[:count+1]
		}
	} else {
		runs, err = sliceRuns(runs, build, count)
		if err != nil {
			return []jenkins.Run{}, err
		}
	}

	if filter == "" {
		return runs, nil
	}

	if validFilter, valid := jenkins.IsValidFilterOption(filter); valid {
		var filteredRuns []jenkins.Run
		for _, run := range runs {
			if run.Status == validFilter {
				filteredRuns = append(filteredRuns, run)
			}
		}

		if len(filteredRuns) == 0 {
			return []jenkins.Run{}, fmt.Errorf("no runs found with status %s in this range", filter)
		}
		return filteredRuns, nil
	} else {
		return []jenkins.Run{}, fmt.Errorf("invalid status filter: %s", filter)
	}
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
	jobCmd.Flags().StringP("build", "b", "lastBuild", "Job build number or (lastBuild, lastFailedBuild, lastCompletedBuild, lastStableBuild, lastUnstableBuild)")
	jobCmd.Flags().StringP("filter-status", "f", "", "Filter job builds by status (success, failure, aborted, unstable, notbuilt)")
	jobCmd.Flags().StringP("stage", "s", "", "Job build stage name")
	jobCmd.Flags().IntP("prev-count", "p", 0, "Number of logs to get preceding selected log")

	rootCmd.AddCommand(jobCmd)
}
