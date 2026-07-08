package analytics

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAnalyticsHandler(c *gin.Context) {

	userID := c.MustGet("userID").(string)

	result, err := GetAnalytics(userID)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, result)
}
