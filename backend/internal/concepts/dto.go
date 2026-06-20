package concepts

type ExtractConceptsRequest struct {
	ResourceID string `json:"resource_id" binding:"required"`
}