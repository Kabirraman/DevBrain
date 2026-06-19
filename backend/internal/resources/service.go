package resources

import (
	"github.com/google/uuid"

	"github.com/Kabirraman/DevBrain/internal/database"
	"github.com/Kabirraman/DevBrain/internal/models"
)

func CreateResource(
	userID string,
	req CreateResourceRequest,
) error {

	uid, err := uuid.Parse(userID)

	if err != nil {
		return err
	}

	resource := models.Resource{
		UserID: uid,

		Title: req.Title,
		Type: req.Type,

		SourceURL: req.SourceURL,
		Content: req.Content,
	}

	return database.DB.Create(&resource).Error
}