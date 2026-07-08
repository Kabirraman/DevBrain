package chat

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ChatHandler(c *gin.Context) {

	var req ChatRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	userID := c.MustGet("userID").(string)

	answer, concepts, err := AskGraph(userID, req.Question)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, ChatResponse{
		Answer:       answer,
		ConceptsUsed: concepts,
	})
}