package gaps

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGapsHandler(c *gin.Context) {

	userID := c.MustGet("userID").(string)

	domain := c.Query("domain")

	if domain == "" {
		domain = "Go"
	}

	result, err := GetGapAnalysis(userID, domain)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, result)
}

func GetDomainsHandler(c *gin.Context) {

	c.JSON(http.StatusOK, DomainsResponse{
		Domains: Domains(),
	})
}
