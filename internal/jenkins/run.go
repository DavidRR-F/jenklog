package jenkins

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

func (r *Run) GetStagesByStatus(status string) []Stage {
	var filteredStages []Stage
	for _, stage := range r.Stages {
		if stage.Status == status {
			filteredStages = append(filteredStages, stage)
		}
	}
	return filteredStages
}
