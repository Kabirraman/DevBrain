package concepts

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Kabirraman/DevBrain/internal/database"
	"github.com/Kabirraman/DevBrain/internal/models"
)

func ExtractRelationshipsHandler(c *gin.Context) {

	var req ExtractConceptsRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	var resource models.Resource

	err := database.DB.
		First(&resource, "id = ?", req.ResourceID).
		Error

	if err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"error": "resource not found",
		})

		return
	}

	relationships, err := ExtractRelationships(
		resource.Content,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	err = SaveRelationships(
		resource.UserID.String(),
		relationships,
	)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"relationships": relationships,
	})
}