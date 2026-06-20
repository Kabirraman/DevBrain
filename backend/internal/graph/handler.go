package graph

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGraphHandler(c *gin.Context) {

	graph, err := BuildGraph()

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(
		http.StatusOK,
		graph,
	)
}