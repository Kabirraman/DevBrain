package resources

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ImportBlogHandler(c *gin.Context) {
	var req ImportBlogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	title, author, description, content, err := ExtractBlogContent(ctx, req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(string)
	if err := CreateResource(userID, CreateResourceRequest{
		Title:       title,
		Type:        "blog",
		Author:      author,
		Description: description,
		SourceURL:   req.URL,
		Content:     content,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "blog imported successfully"})
}
