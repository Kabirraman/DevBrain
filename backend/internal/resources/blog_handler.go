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

	userID := c.MustGet("userID").(string)

resourceReq := CreateResourceRequest{
	Title:     title,
	Type:      "blog",
	SourceURL: req.URL,
	Content:   content,
}

err = CreateResource(
	userID,
	resourceReq,
)

if err != nil {

	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})

	return
}

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
	"message": "blog imported successfully",
})
}