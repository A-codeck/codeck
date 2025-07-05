package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int            `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string         `gorm:"type:text;unique;not null" json:"email"`
	Name      string         `gorm:"type:text;not null" json:"name"`
	Password  string         `gorm:"type:text;not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}
