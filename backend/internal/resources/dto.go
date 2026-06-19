package resources

type CreateResourceRequest struct {
	Title     string `json:"title" binding:"required"`
	Type      string `json:"type" binding:"required"`
	SourceURL string `json:"source_url"`
	Content   string `json:"content" binding:"required"`
}