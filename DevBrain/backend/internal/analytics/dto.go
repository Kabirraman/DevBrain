package analytics

type CoverageEntry struct {
	Domain     string `json:"domain"`
	Completion int    `json:"completion"`
}

type AnalyticsResponse struct {
	TotalResources    int             `json:"total_resources"`
	TotalConcepts     int             `json:"total_concepts"`
	MostLearnedTopics []string        `json:"most_learned_topics"`
	KnowledgeCoverage []CoverageEntry `json:"knowledge_coverage"`
	MissingConcepts   []string        `json:"missing_concepts"`
}
