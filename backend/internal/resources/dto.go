package resources

type CreateResourceRequest struct {
	Title       string `json:"title" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Author      string `json:"author,omitempty"`
	Description string `json:"description,omitempty"`
	SourceURL   string `json:"source_url"`
	Content     string `json:"content" binding:"required"`
}
