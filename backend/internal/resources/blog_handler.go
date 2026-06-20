package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ImportBlogHandler(c *gin.Context) {

	var req ImportBlogRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	title, content, err := ExtractBlogContent(
		req.URL,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"title":   title,
		"content": content[:300],
	})
}