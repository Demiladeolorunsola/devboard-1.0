package services

import (
	"errors"
	"time"

	"github.com/Demiladeolorunsola/devboard-1.0/internal/models"
	"gorm.io/gorm"
)

type TimeService struct {
	DB *gorm.DB
}

func NewTimeService(db *gorm.DB) *TimeService {
	return &TimeService{
		DB: db,
	}
}

// =========================
// START TIMER
// =========================

func (s *TimeService) Start(
	userID uint,
	projectID uint,
	taskID uint,
) (*models.TimeEntry, error) {

	// Make sure the project belongs to the user.
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

	// Make sure the task belongs to the user and project.
	var task models.Task

	result = s.DB.
		Where(
			"id = ? AND project_id = ? AND user_id = ?",
			taskID,
			projectID,
			userID,
		).
		First(&task)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("task not found")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	// Prevent multiple active timers for the same user.
	var activeTimer models.TimeEntry

	result = s.DB.
		Where(
			"user_id = ? AND end_time IS NULL",
			userID,
		).
		First(&activeTimer)

	if result.Error == nil {
		return nil, errors.New("you already have an active timer")
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	timeEntry := models.TimeEntry{
		UserID:    userID,
		ProjectID: projectID,
		TaskID:    taskID,
		StartTime: time.Now(),
		Duration:  0,
	}

	if err := s.DB.Create(&timeEntry).Error; err != nil {
		return nil, err
	}

	return &timeEntry, nil
}

// =========================
// STOP TIMER
// =========================

func (s *TimeService) Stop(
	userID uint,
	timeEntryID uint,
) (*models.TimeEntry, error) {

	var timeEntry models.TimeEntry

	result := s.DB.
		Where(
			"id = ? AND user_id = ?",
			timeEntryID,
			userID,
		).
		First(&timeEntry)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("time entry not found")
	}

	if result.Error != nil {
		return nil, result.Error
	}

	if timeEntry.EndTime != nil {
		return nil, errors.New("timer has already been stopped")
	}

	endTime := time.Now()

	duration := int64(endTime.Sub(timeEntry.StartTime).Seconds())

	if duration < 0 {
		duration = 0
	}

	timeEntry.EndTime = &endTime
	timeEntry.Duration = duration

	if err := s.DB.Save(&timeEntry).Error; err != nil {
		return nil, err
	}

	return &timeEntry, nil
}

// =========================
// GET USER TIME ENTRIES
// =========================

func (s *TimeService) GetAll(
	userID uint,
) ([]models.TimeEntry, error) {

	var entries []models.TimeEntry

	if err := s.DB.
		Where("user_id = ?", userID).
		Order("start_time DESC").
		Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

// =========================
// GET TIME ENTRIES FOR PROJECT
// =========================

func (s *TimeService) GetByProject(
	userID uint,
	projectID uint,
) ([]models.TimeEntry, error) {

	// Make sure the project belongs to the user.
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

	var entries []models.TimeEntry

	if err := s.DB.
		Where(
			"user_id = ? AND project_id = ?",
			userID,
			projectID,
		).
		Order("start_time DESC").
		Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}

// =========================
// GET TIME ENTRIES FOR TASK
// =========================

func (s *TimeService) GetByTask(
	userID uint,
	taskID uint,
) ([]models.TimeEntry, error) {

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

	var entries []models.TimeEntry

	if err := s.DB.
		Where(
			"user_id = ? AND task_id = ?",
			userID,
			taskID,
		).
		Order("start_time DESC").
		Find(&entries).Error; err != nil {
		return nil, err
	}

	return entries, nil
}
