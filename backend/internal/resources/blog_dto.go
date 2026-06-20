package resources

type ImportBlogRequest struct {
	URL string `json:"url" binding:"required,url"`
}