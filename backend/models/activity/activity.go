package activity

import (
	"time"

	"gorm.io/gorm"
)

type Activity struct {
	ID            int            `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatorID     int            `gorm:"not null;index" json:"creator_id"`
	Title         string         `gorm:"type:text;not null" json:"title"`
	Date          time.Time      `gorm:"type:date;not null" json:"date"`
	ActivityImage string         `gorm:"type:text;not null" json:"activity_image"`
	Description   *string        `gorm:"type:text" json:"description,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`

	// Many-to-many relationship with groups
	Groups []int `gorm:"-" json:"group_ids,omitempty"` // This will be populated from ActivityGroups
}

// ActivityWithGroups represents an activity with its associated group information
type ActivityWithGroups struct {
	Activity
	GroupNames []string `json:"group_names,omitempty"`
}
