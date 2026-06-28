package chat

type ChatRequest struct {
	Question string `json:"question"`
}

type ChatResponse struct {
	Answer      string   `json:"answer"`
	ConceptsUsed []string `json:"concepts_used"`
}