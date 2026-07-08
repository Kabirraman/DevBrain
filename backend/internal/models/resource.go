package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Resource struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`

	UserID    uuid.UUID `gorm:"type:uuid;not null"`

	Title     string
	Type      string

	SourceURL string
	Content   string `gorm:"type:text"`

	ConceptsExtracted bool `gorm:"default:false"`

	CreatedAt time.Time
}

func (r *Resource) BeforeCreate(tx *gorm.DB) error {
	r.ID = uuid.New()
	return nil
}