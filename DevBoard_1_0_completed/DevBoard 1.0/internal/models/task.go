package models

import "time"

type Task struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"not null" json:"title"`
	Description string `json:"description"`

	Status   string `gorm:"not null;default:todo" json:"status"`
	Priority string `gorm:"not null;default:medium" json:"priority"`

	UserID    uint `gorm:"not null;index" json:"user_id"`
	ProjectID uint `gorm:"not null;index" json:"project_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User    User    `gorm:"foreignKey:UserID" json:"-"`
	Project Project `gorm:"foreignKey:ProjectID" json:"-"`
}
