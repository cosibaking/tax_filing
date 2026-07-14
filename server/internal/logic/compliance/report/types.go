package report

const (
	StructuredSchemaVersion = 1
	ConclusionNormal        = "normal"
	ConclusionAttention     = "attention"
	ConclusionUrgent        = "urgent"
)

type StructuredReport struct {
	SchemaVersion int              `json:"schemaVersion"`
	Summary       ReportSummary    `json:"summary"`
	Categories    []ReportCategory `json:"categories"`
	Anomalies     []ReportAnomaly  `json:"anomalies"`
}

type ReportSummary struct {
	Conclusion       string  `json:"conclusion"`
	CompletenessRate float64 `json:"completenessRate"`
	HighCount        int     `json:"highCount"`
	MediumCount      int     `json:"mediumCount"`
	LowCount         int     `json:"lowCount"`
	DataNotice       string  `json:"dataNotice"`
}

type ReportCategory struct {
	Code    string        `json:"code"`
	Name    string        `json:"name"`
	Status  string        `json:"status"`
	Summary string        `json:"summary"`
	Checks  []ReportCheck `json:"checks"`
}

type ReportCheck struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ReportAnomaly struct {
	Code                 string   `json:"code"`
	CategoryCode         string   `json:"categoryCode"`
	Title                string   `json:"title"`
	Severity             string   `json:"severity"`
	Facts                string   `json:"facts"`
	Basis                string   `json:"basis"`
	Impact               string   `json:"impact"`
	Recommendation       string   `json:"recommendation"`
	RequiredMaterials    []string `json:"requiredMaterials"`
	DueDate              string   `json:"dueDate"`
	RequiresManualReview bool     `json:"requiresManualReview"`
	RuleVersion          string   `json:"ruleVersion"`
}
