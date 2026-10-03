package services

import (
	"errors"
	"strings"

	"github.com/Demiladeolorunsola/devboard-1.0/internal/models"
	"gorm.io/gorm"
)

type ProjectService struct {
	DB *gorm.DB
}

func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{
		DB: db,
	}
}

func (s *ProjectService) Create(
	userID uint,
	name string,
	description string,
) (*models.Project, error) {

	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errors.New("project name is required")
	}

	project := models.Project{
		Name:        name,
		Description: strings.TrimSpace(description),
		Status:      "active",
		UserID:      userID,
	}

	if err := s.DB.Create(&project).Error; err != nil {
		return nil, err
	}

	return &project, nil
}

func (s *ProjectService) GetAll(
	userID uint,
) ([]models.Project, error) {

	var projects []models.Project

	if err := s.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *ProjectService) GetByID(
	userID uint,
	projectID uint,
) (*models.Project, error) {

	var project models.Project

	result := s.DB.
		Where("id = ? AND user_id = ?", projectID, userID).
		First(&project)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("project not found")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &project, nil
}

func (s *ProjectService) Update(
	userID uint,
	projectID uint,
	name string,
	description string,
	status string,
) (*models.Project, error) {

	project, err := s.GetByID(userID, projectID)

	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(name) != "" {
		project.Name = strings.TrimSpace(name)
	}

	project.Description = strings.TrimSpace(description)

	if status != "" {
		allowedStatuses := map[string]bool{
			"active":    true,
			"completed": true,
			"archived":  true,
		}

		if !allowedStatuses[status] {
			return nil, errors.New("invalid project status")
		}

		project.Status = status
	}

	if err := s.DB.Save(project).Error; err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) Delete(
	userID uint,
	projectID uint,
) error {

	project, err := s.GetByID(userID, projectID)

	if err != nil {
		return err
	}

	if err := s.DB.Delete(project).Error; err != nil {
		return err
	}

	return nil
}
