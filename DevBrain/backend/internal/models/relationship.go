package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Relationship struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index"`

	Source   string
	Relation string
	Target   string

	CreatedAt time.Time
}

func (r *Relationship) BeforeCreate(tx *gorm.DB) error {
	r.ID = uuid.New()
	return nil
}