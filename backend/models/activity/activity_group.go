package activity

import (
	"time"

	"gorm.io/gorm"
)

// ActivityGroup represents the many-to-many relationship between activities and groups
type ActivityGroup struct {
	ID         int            `gorm:"primaryKey;autoIncrement" json:"id"`
	ActivityID int            `gorm:"not null;index" json:"activity_id"`
	GroupID    int            `gorm:"not null;index" json:"group_id"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}
