package jenkins

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/DavidRR-F/jenklog/internal/config"
)

const (
	LAST_BUILD            = "lastBuild"
	LAST_FAILED_BUILD     = "lastFailedBuild"
	LAST_SUCCESSFUL_BUILD = "lastCompletedBuild"
	LAST_STABLE_BUILD     = "lastStableBuild"
	LAST_UNSTABLE_BUILD   = "lastUnstableBuild"
)

const (
	JOB_RUNS = "/job/%s/wfapi/runs"
	JOB_LOGS = "/job/%s/%s/consoleText"
)

func GetJenkinsJobRuns(job string) ([]Run, error) {
	config, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf(JOB_RUNS, job)
	body, err := queryJenkins(config, endpoint)

	if err != nil {
		return nil, err
	}

	var runs []Run
	if err := json.Unmarshal(body, &runs); err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %v", err)
	}
	return runs, nil
}

func GetJenkinsJobLog(job, build string) (Log, error) {
	config, err := config.GetConfig()
	if err != nil {
		return Log{}, err
	}

	endpoint := fmt.Sprintf(JOB_LOGS, job, build)
	body, err := queryJenkins(config, endpoint)

	if err != nil {
		return Log{}, err
	}

	return Log{bytes: body}, nil
}

func queryJenkins(config *config.Config, endpoint string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", config.URL+endpoint, nil)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(config.Username, config.Token)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
