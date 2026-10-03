package models

import "time"

type TimeEntry struct {
	ID uint `gorm:"primaryKey" json:"id"`

	UserID    uint `gorm:"not null;index" json:"user_id"`
	ProjectID uint `gorm:"not null;index" json:"project_id"`
	TaskID    uint `gorm:"not null;index" json:"task_id"`

	StartTime time.Time  `gorm:"not null" json:"start_time"`
	EndTime   *time.Time `json:"end_time"`

	Duration int64 `gorm:"not null;default:0" json:"duration"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User    User    `gorm:"foreignKey:UserID" json:"-"`
	Project Project `gorm:"foreignKey:ProjectID" json:"-"`
	Task    Task    `gorm:"foreignKey:TaskID" json:"-"`
}
