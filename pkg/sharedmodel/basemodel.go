package sharedmodel

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel for standard fields
type BaseModel struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();index" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
