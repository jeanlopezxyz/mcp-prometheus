package prometheus

// APIResponse is the common wrapper for all Prometheus API responses.
type APIResponse struct {
	Status    string      `json:"status"`
	Data      interface{} `json:"data,omitempty"`
	ErrorType string      `json:"errorType,omitempty"`
	Error     string      `json:"error,omitempty"`
}

// QueryData holds the result of an instant or range query.
type QueryData struct {
	ResultType string        `json:"resultType"`
	Result     []interface{} `json:"result"`
}

// TargetsData holds scrape target information.
type TargetsData struct {
	ActiveTargets  []ActiveTarget  `json:"activeTargets"`
	DroppedTargets []DroppedTarget `json:"droppedTargets"`
}

// ActiveTarget represents an active scrape target.
type ActiveTarget struct {
	DiscoveredLabels map[string]string `json:"discoveredLabels"`
	Labels           map[string]string `json:"labels"`
	ScrapePool       string            `json:"scrapePool"`
	ScrapeURL        string            `json:"scrapeUrl"`
	GlobalURL        string            `json:"globalUrl"`
	LastError        string            `json:"lastError"`
	LastScrape       string            `json:"lastScrape"`
	LastScrapeDuration float64         `json:"lastScrapeDuration"`
	Health           string            `json:"health"`
	ScrapeInterval   string            `json:"scrapeInterval"`
	ScrapeTimeout    string            `json:"scrapeTimeout"`
}

// DroppedTarget represents a dropped scrape target.
type DroppedTarget struct {
	DiscoveredLabels map[string]string `json:"discoveredLabels"`
}

// RulesData holds alerting and recording rules.
type RulesData struct {
	Groups []RuleGroup `json:"groups"`
}

// RuleGroup represents a group of rules.
type RuleGroup struct {
	Name     string `json:"name"`
	File     string `json:"file"`
	Rules    []Rule `json:"rules"`
	Interval float64 `json:"interval"`
}

// Rule represents an alerting or recording rule.
type Rule struct {
	State       string            `json:"state,omitempty"`
	Name        string            `json:"name"`
	Query       string            `json:"query"`
	Duration    float64           `json:"duration,omitempty"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Alerts      []Alert           `json:"alerts,omitempty"`
	Health      string            `json:"health"`
	LastError   string            `json:"lastError,omitempty"`
	Type        string            `json:"type"`
}

// Alert represents an alert instance within a rule.
type Alert struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	State       string            `json:"state"`
	ActiveAt    string            `json:"activeAt"`
	Value       string            `json:"value"`
}

// BuildInfo holds Prometheus build information.
type BuildInfo struct {
	Version   string `json:"version"`
	Revision  string `json:"revision"`
	Branch    string `json:"branch"`
	BuildUser string `json:"buildUser"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
}

// RuntimeInfo holds Prometheus runtime information.
type RuntimeInfo struct {
	StartTime           string `json:"startTime"`
	CWD                 string `json:"CWD"`
	ReloadConfigSuccess bool   `json:"reloadConfigSuccess"`
	LastConfigTime      string `json:"lastConfigTime"`
	CorruptionCount     int    `json:"corruptionCount"`
	GoroutineCount      int    `json:"goroutineCount"`
	GOMAXPROCS          int    `json:"GOMAXPROCS"`
	GOGC                string `json:"GOGC"`
	GODEBUG             string `json:"GODEBUG"`
	StorageRetention    string `json:"storageRetention"`
}
