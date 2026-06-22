package concepts

import (
	"github.com/Kabirraman/DevBrain/internal/database"
	"github.com/Kabirraman/DevBrain/internal/models"
	"net/http"
	"github.com/gin-gonic/gin"
)


type ConceptDetailResponse struct {
	Concept       string                    `json:"concept"`
	Relationships []RelationshipDTO         `json:"relationships"`
}

func GetConceptDetails(
	name string,
) (*ConceptDetailResponse, error) {

	var relationships []models.Relationship

	err := database.DB.
		Where(
			"source = ? OR target = ?",
			name,
			name,
		).
		Find(&relationships).
		Error

	if err != nil {
		return nil, err
	}

	result := []RelationshipDTO{}

	for _, rel := range relationships {

		result = append(
			result,
			RelationshipDTO{
				Source: rel.Source,
				Relation: rel.Relation,
				Target: rel.Target,
			},
		)
	}

	return &ConceptDetailResponse{
		Concept: name,
		Relationships: result,
	}, nil
}
func GetConceptDetailsHandler(
	c *gin.Context,
) {

	name := c.Param("name")

	result, err :=
		GetConceptDetails(name)

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		result,
	)
}