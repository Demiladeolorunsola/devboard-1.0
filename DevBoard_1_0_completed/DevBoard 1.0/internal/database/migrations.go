package database

import (
	"fmt"

	"github.com/Demiladeolorunsola/devboard-1.0/internal/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Task{},
		&models.TimeEntry{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}
