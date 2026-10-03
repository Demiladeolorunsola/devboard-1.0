package services

import (
	"errors"
	"strings"

	"github.com/Demiladeolorunsola/devboard-1.0/internal/models"
	"gorm.io/gorm"
)

type TaskService struct {
	DB *gorm.DB
}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{
		DB: db,
	}
}

// Create creates a task inside a project owned by the user.
func (s *TaskService) Create(
	userID uint,
	projectID uint,
	title string,
	description string,
	priority string,
) (*models.Task, error) {

	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)
	priority = strings.TrimSpace(priority)

	if title == "" {
		return nil, errors.New("task title is required")
	}

	if priority == "" {
		priority = "medium"
	}

	allowedPriorities := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
	}

	if !allowedPriorities[priority] {
		return nil, errors.New("invalid task priority")
	}

	// Make sure the project belongs to the authenticated user.
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

	task := models.Task{
		Title:       title,
		Description: description,
		Status:      "todo",
		Priority:    priority,
		UserID:      userID,
		ProjectID:   projectID,
	}

	if err := s.DB.Create(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

// GetByProject returns all tasks belonging to a project owned by the user.
func (s *TaskService) GetByProject(
	userID uint,
	projectID uint,
) ([]models.Task, error) {

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

	var tasks []models.Task

	if err := s.DB.
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Order("created_at DESC").
		Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetByID returns one task if it belongs to the authenticated user.
func (s *TaskService) GetByID(
	userID uint,
	taskID uint,
) (*models.Task, error) {

	var task models.Task

	result := s.DB.
		Where("id = ? AND user_id = ?", taskID, userID).
		First(&task)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("task not found")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return &task, nil
}

// Update updates a task owned by the authenticated user.
func (s *TaskService) Update(
	userID uint,
	taskID uint,
	title string,
	description string,
	status string,
	priority string,
) (*models.Task, error) {

	task, err := s.GetByID(userID, taskID)

	if err != nil {
		return nil, err
	}

	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)
	status = strings.TrimSpace(status)
	priority = strings.TrimSpace(priority)

	if title != "" {
		task.Title = title
	}

	task.Description = description

	if status != "" {
		allowedStatuses := map[string]bool{
			"todo":        true,
			"in_progress": true,
			"completed":   true,
		}

		if !allowedStatuses[status] {
			return nil, errors.New("invalid task status")
		}

		task.Status = status
	}

	if priority != "" {
		allowedPriorities := map[string]bool{
			"low":    true,
			"medium": true,
			"high":   true,
		}

		if !allowedPriorities[priority] {
			return nil, errors.New("invalid task priority")
		}

		task.Priority = priority
	}

	if err := s.DB.Save(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

// Delete deletes a task owned by the authenticated user.
func (s *TaskService) Delete(
	userID uint,
	taskID uint,
) error {

	task, err := s.GetByID(userID, taskID)

	if err != nil {
		return err
	}

	if err := s.DB.Delete(task).Error; err != nil {
		return err
	}

	return nil
}
