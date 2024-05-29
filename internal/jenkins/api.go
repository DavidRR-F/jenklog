package jenkins

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/DavidRR-F/jenklog/internal/config"
)

const (
	JOB_LOGS = "/job/%s/%s/consoleText"
	JOB_INFO = "/job/%s/%s/api/json"
)

const (
	LAST_BUILD            = "lastBuild"
	LAST_FAILED_BUILD     = "lastFailedBuild"
	LAST_SUCCESSFUL_BUILD = "lastCompletedBuild"
	LAST_STABLE_BUILD     = "lastStableBuild"
	LAST_UNSTABLE_BUILD   = "lastUnstableBuild"
)

var validBuildOptions = []string{
	LAST_BUILD,
	LAST_FAILED_BUILD,
	LAST_SUCCESSFUL_BUILD,
	LAST_STABLE_BUILD,
	LAST_UNSTABLE_BUILD,
}

func IsValidBuildOption(option string) bool {
	for _, validOption := range validBuildOptions {
		if option == validOption {
			return true
		}
	}

	if _, err := strconv.Atoi(option); err == nil {
		return true
	}

	return false
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

	return Log{id: build, stage: "All", bytes: body}, nil
}

func GetJenkinsJobInfo(job, build string) (BuildInfo, error) {
	config, err := config.GetConfig()
	if err != nil {
		return BuildInfo{}, err
	}
	endpoint := fmt.Sprintf(JOB_INFO, job, build)
	body, err := queryJenkins(config, endpoint)

	if err != nil {
		return BuildInfo{}, err
	}

	var buildInfo BuildInfo
	if err := json.Unmarshal(body, &buildInfo); err != nil {
		return BuildInfo{}, fmt.Errorf("error unmarshalling response body: %v", err)
	}
	return buildInfo, nil

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
