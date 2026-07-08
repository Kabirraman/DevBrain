package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserConcept struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_user_concept,unique"`
	ConceptID uuid.UUID `gorm:"type:uuid;not null;index:idx_user_concept,unique"`

	CreatedAt time.Time
}

func (uc *UserConcept) BeforeCreate(tx *gorm.DB) error {
	uc.ID = uuid.New()
	return nil
}
