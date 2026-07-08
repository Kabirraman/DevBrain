package search

type SearchRequest struct {
	Query string `json:"query" binding:"required"`
}

type SearchResult struct {
	ResourceID string `json:"resource_id"`
	Title      string `json:"title"`
	Snippet    string `json:"snippet"`
	Score      int    `json:"score"`
}

type SearchResponse struct {
	Answer  string         `json:"answer"`
	Results []SearchResult `json:"results"`
}
