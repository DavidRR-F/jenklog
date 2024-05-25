package jenkins

type BuildInfo struct {
	Class             string         `json:"_class"`
	Actions           []Action       `json:"actions"`
	Artifacts         []interface{}  `json:"artifacts"`
	Building          bool           `json:"building"`
	Description       *string        `json:"description"`
	DisplayName       string         `json:"displayName"`
	Duration          int            `json:"duration"`
	EstimatedDuration int            `json:"estimatedDuration"`
	Executor          *interface{}   `json:"executor"`
	FullDisplayName   string         `json:"fullDisplayName"`
	ID                string         `json:"id"`
	KeepLog           bool           `json:"keepLog"`
	Number            int            `json:"number"`
	QueueID           int            `json:"queueId"`
	Result            string         `json:"result"`
	Timestamp         int64          `json:"timestamp"`
	URL               string         `json:"url"`
	ChangeSets        []interface{}  `json:"changeSets"`
	Culprits          []interface{}  `json:"culprits"`
	InProgress        bool           `json:"inProgress"`
	NextBuild         *interface{}   `json:"nextBuild"`
	PreviousBuild     *PreviousBuild `json:"previousBuild"`
}

type Action struct {
	Class  string  `json:"_class,omitempty"`
	Causes []Cause `json:"causes,omitempty"`
}

type Cause struct {
	Class            string  `json:"_class"`
	ShortDescription string  `json:"shortDescription"`
	UserID           *string `json:"userId"`
	UserName         string  `json:"userName"`
}

type PreviousBuild struct {
	Number int    `json:"number"`
	URL    string `json:"url"`
}
