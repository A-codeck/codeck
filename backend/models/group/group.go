package group

import (
	"time"

	"gorm.io/gorm"
)

type Group struct {
	ID          int            `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatorID   int            `gorm:"not null;index" json:"creator_id"`
	Name        string         `gorm:"type:text;not null" json:"name"`
	StartDate   time.Time      `gorm:"type:date;not null" json:"start_date"`
	EndDate     time.Time      `gorm:"type:date;not null" json:"end_date"`
	GroupImage  *string        `gorm:"type:text" json:"group_image,omitempty"`
	Description *string        `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}

type GroupMember struct {
	UserID   int     `gorm:"not null;index" json:"user_id"`
	GroupID  int     `gorm:"not null;index" json:"group_id"`
	Nickname *string `gorm:"type:text" json:"nickname,omitempty"`
}

type GroupInvite struct {
	InviteCode string         `gorm:"primaryKey;type:text" json:"invite_code"`
	GroupID    int            `gorm:"not null;index" json:"group_id"`
	CreatedBy  int            `gorm:"not null;index" json:"created_by"`
	CreatedAt  time.Time      `json:"created_at"`
	ExpiresAt  *time.Time     `json:"expires_at,omitempty"`
	IsActive   bool           `gorm:"default:true" json:"is_active"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}
