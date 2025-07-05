package activity

import (
	"time"

	"gorm.io/gorm"
)

type Activity struct {
	ID            int            `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatorID     int            `gorm:"not null;index" json:"creator_id"`
	GroupID       int            `gorm:"not null;index" json:"group_id"`
	Title         string         `gorm:"type:text;not null" json:"title"`
	Date          time.Time      `gorm:"type:date;not null" json:"date"`
	ActivityImage *string        `gorm:"type:text" json:"activity_image,omitempty"`
	Description   *string        `gorm:"type:text" json:"description,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}
