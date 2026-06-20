package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Concept struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`

	Name      string `gorm:"not null"`
	Category  string

	CreatedAt time.Time
}

func (c *Concept) BeforeCreate(tx *gorm.DB) error {
	c.ID = uuid.New()
	return nil
}