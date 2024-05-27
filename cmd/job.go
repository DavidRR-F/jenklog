/*
2024 David Rose-Franklin <david.franklin.dev@gmail.ocm>
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
	Run:   executeJob,
}

func executeJob(cmd *cobra.Command, args []string) {
	job := args[0]

	stage, err := cmd.Flags().GetString("stage")
	if err != nil {
		fmt.Println("Error getting stage flag:", err)
		os.Exit(1)
	}

	build, err := cmd.Flags().GetString("build")
	if err != nil {
		fmt.Println("Error getting build flag:", err)
		os.Exit(1)
	}

	count, err := cmd.Flags().GetInt("prev-count")
	if err != nil {
		fmt.Println("Error getting prev-count flag:", err)
		os.Exit(1)
	}

	filter, err := cmd.Flags().GetString("filter-status")
	if err != nil {
		fmt.Println("Error getting filter-status flag:", err)
		os.Exit(1)
	}

	if !jenkins.IsValidBuildOption(build) {
		fmt.Println("Not a valid build option")
		os.Exit(1)
	}

	switch count {
	case 0:
		err := handleSingleLogRequest(job, stage, build, filter)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		err := handleMultiLogRequest(job, stage, build, count, filter)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func init() {
	jobCmd.Flags().StringP("build", "b", "lastBuild", "Job build number or (lastBuild, lastFailedBuild, lastCompletedBuild, lastStableBuild, lastUnstableBuild)")
	jobCmd.Flags().StringP("filter-status", "f", "", "Filter job builds by status (success, failure, aborted, unstable, notbuilt)")
	jobCmd.Flags().StringP("stage", "s", "", "Job build stage name")
	jobCmd.Flags().IntP("prev-count", "p", 0, "Number of logs to get preceding selected log")

	rootCmd.AddCommand(jobCmd)
}

func handleSingleLogRequest(job, stage, build, filter string) error {
	if filter != "" {
		return fmt.Errorf("must select multiple builds to fitler by status")
	}
	log, err := getJobLog(job, stage, build)
	if err != nil {
		return err
	}
	log.Print()
	return nil
}

func handleMultiLogRequest(job, stage, build string, count int, filter string) error {
	runs, err := getJobRuns(job, build, filter, count)
	if err != nil {
		return err
	}
	logs := make([]jenkins.Log, len(runs))
	var wg sync.WaitGroup
	for i, run := range runs {
		wg.Add(1)
		go getJobLogs(job, stage, run, &wg, &logs[i])
	}
	wg.Wait()

	for _, log := range logs {
		log.Print()
	}
	return nil
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

func getJobLogs(job, stage string, run jenkins.Run, wg *sync.WaitGroup, logDest *jenkins.Log) {
	defer wg.Done()
	log, err := getJobLog(job, stage, run.ID)
	if err != nil {
		fmt.Println(err)
		return
	}
	*logDest = log
}

func getJobRuns(job, build, filter string, count int) (jenkins.Runs, error) {
	if _, err := strconv.Atoi(build); err != nil && build != jenkins.LAST_BUILD {
		buildInfo, err := jenkins.GetJenkinsJobInfo(job, build)
		if err != nil {
			return jenkins.Runs{}, err
		}
		build = buildInfo.ID
	}

	runs, err := jenkins.GetJenkinsJobRuns(job)

	if err != nil {
		return jenkins.Runs{}, err
	}

	if build == jenkins.LAST_BUILD {
		if count+1 < len(runs) {
			runs = runs[:count+1]
		}
	} else {
		runs, err = runs.Slice(build, count)
		if err != nil {
			return jenkins.Runs{}, err
		}
	}

	if filter == "" {
		return runs, nil
	}

	if validFilter, valid := jenkins.IsValidFilterOption(filter); valid {
		return runs.Filter(validFilter)
	} else {
		return jenkins.Runs{}, fmt.Errorf("invalid status filter: %s", filter)
	}
}
