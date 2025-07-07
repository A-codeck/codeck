package comment

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID         int            `gorm:"primaryKey;autoIncrement" json:"id"`
	ActivityID int            `gorm:"not null;index" json:"activity_id"`
	UserID     int            `gorm:"not null;index" json:"user_id"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
