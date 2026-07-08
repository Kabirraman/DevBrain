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

func GetResources(userID string) ([]models.Resource, error) {

	var resources []models.Resource

	err := database.DB.
		Where("user_id = ?", userID).
		Find(&resources).
		Error

	if err != nil {
		return nil, err
	}

	return resources, nil
}

func GetResourceByID(
	userID string,
	resourceID string,
) (*models.Resource, error) {

	var resource models.Resource

	err := database.DB.
		Where(
			"id = ? AND user_id = ?",
			resourceID,
			userID,
		).
		First(&resource).
		Error

	if err != nil {
		return nil, err
	}

	return &resource, nil
}