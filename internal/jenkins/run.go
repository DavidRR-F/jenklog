package jenkins

import "fmt"

type Run struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	StartTime int64   `json:"startTimeMillis"`
	EndTime   int64   `json:"endTimeMillis"`
	Duration  int64   `json:"durationMillis"`
	Stages    []Stage `json:"stages"`
	Log       Log
}

type Stage struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	StartTime int64  `json:"startTimeMillis"`
	Duration  int64  `json:"durationMillis"`
	Links     Links  `json:"_links"`
}

type Links struct {
	Self Link `json:"self"`
}

type Link struct {
	Href string `json:"href"`
}

type Runs []Run

func (r Runs) Filter(status string) (Runs, error) {
	var filtered Runs
	for _, run := range r {
		if run.Status == status {
			filtered = append(filtered, run)
		}
	}
	if len(filtered) == 0 {
		return nil, fmt.Errorf("no runs found with status %s in this range", status)
	}
	return filtered, nil
}

func (r Runs) Slice(build string, count int) (Runs, error) {
	startIndex := -1
	for i, run := range r {
		if run.ID == build {
			startIndex = i
			break
		}
	}

	if startIndex == -1 {
		return nil, fmt.Errorf("value %s not found in the builds", build)
	}

	endIndex := startIndex + count + 1
	if endIndex > len(r) {
		endIndex = len(r)
	}

	return r[startIndex:endIndex], nil
}
