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

	if !jenkins.IsValidBuildOption(build) {
		fmt.Println("Not a valid build option")
		os.Exit(1)
	}

	switch count {
	case 0:
		err := handleSingleLogRequest(job, stage, build)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		err := handleMultiLogRequest(job, stage, build, count)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func init() {
	jobCmd.Flags().StringP("build", "b", "lastBuild", "Job build number or (lastBuild, lastFailedBuild, lastCompletedBuild, lastStableBuild, lastUnstableBuild)")
	jobCmd.Flags().StringP("stage", "s", "", "Job build stage name")
	jobCmd.Flags().IntP("prev-count", "p", 0, "Number of logs to get preceding selected log")

	rootCmd.AddCommand(jobCmd)
}

func handleSingleLogRequest(job, stage, build string) error {
	log, err := getJobLog(job, stage, build)
	if err != nil {
		return err
	}
	log.Print()
	return nil
}

func handleMultiLogRequest(job, stage, build string, count int) error {
	if _, err := strconv.Atoi(build); err != nil {
		buildInfo, err := jenkins.GetJenkinsJobInfo(job, build)
		if err != nil {
			return err
		}
		build = buildInfo.ID
	}

	buildInt, err := strconv.Atoi(build)

	if err != nil {
		return err
	}

	if buildInt-count < 1 {
		count = buildInt
	}

	var ids []string
	for i := buildInt; i > buildInt-count; i-- {
		ids = append(ids, strconv.Itoa(i))
	}

	logs := make([]jenkins.Log, len(ids))
	var wg sync.WaitGroup
	for i, id := range ids {
		wg.Add(1)
		go getJobLogs(job, stage, id, &wg, &logs[i])
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

func getJobLogs(job, stage string, id string, wg *sync.WaitGroup, logDest *jenkins.Log) {
	defer wg.Done()
	log, err := getJobLog(job, stage, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	*logDest = log
}
