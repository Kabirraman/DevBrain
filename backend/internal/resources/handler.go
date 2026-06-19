package resources

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateResourceHandler(c *gin.Context) {

	var req CreateResourceRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	userID := c.MustGet("userID").(string)

	err := CreateResource(
		userID,
		req,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "resource created",
	})
}