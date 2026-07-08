package search

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchHandler(c *gin.Context) {

	var req SearchRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	userID := c.MustGet("userID").(string)

	result, err := Search(userID, req.Query)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, result)
}
