package services

import (
	"errors"

	"github.com/Demiladeolorunsola/devboard-1.0/internal/models"
	"gorm.io/gorm"
)

type ReportService struct {
	DB *gorm.DB
}

func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{
		DB: db,
	}
}

// =========================
// PROJECT TIME SUMMARY
// =========================

type ProjectTimeSummary struct {
	ProjectID      uint    `json:"project_id"`
	ProjectName    string  `json:"project_name"`
	TotalSeconds   int64   `json:"total_seconds"`
	TotalHours     float64 `json:"total_hours"`
	EntryCount     int     `json:"entry_count"`
	HasActiveTimer bool    `json:"has_active_timer"`
}

// ProjectTimeSummary returns how much time has been logged against a
// single project owned by the authenticated user.
func (s *ReportService) ProjectTimeSummary(
	userID uint,
	projectID uint,
) (*ProjectTimeSummary, error) {

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
		Where("user_id = ? AND project_id = ?", userID, projectID).
		Find(&entries).Error; err != nil {
		return nil, err
	}

	var totalSeconds int64
	hasActive := false

	for _, entry := range entries {
		totalSeconds += entry.Duration

		if entry.EndTime == nil {
			hasActive = true
		}
	}

	return &ProjectTimeSummary{
		ProjectID:      project.ID,
		ProjectName:    project.Name,
		TotalSeconds:   totalSeconds,
		TotalHours:     roundToTwoDecimals(float64(totalSeconds) / 3600),
		EntryCount:     len(entries),
		HasActiveTimer: hasActive,
	}, nil
}

// =========================
// TASK PRODUCTIVITY
// =========================

type TaskProductivityEntry struct {
	TaskID       uint    `json:"task_id"`
	Title        string  `json:"title"`
	Status       string  `json:"status"`
	Priority     string  `json:"priority"`
	TotalSeconds int64   `json:"total_seconds"`
	TotalHours   float64 `json:"total_hours"`
	EntryCount   int     `json:"entry_count"`
}

// TaskProductivity returns time-tracking totals for every task owned by
// the authenticated user, ordered by time spent (most first).
func (s *ReportService) TaskProductivity(
	userID uint,
) ([]TaskProductivityEntry, error) {

	var tasks []models.Task

	if err := s.DB.
		Where("user_id = ?", userID).
		Find(&tasks).Error; err != nil {
		return nil, err
	}

	var entries []models.TimeEntry

	if err := s.DB.
		Where("user_id = ?", userID).
		Find(&entries).Error; err != nil {
		return nil, err
	}

	totalsByTask := make(map[uint]int64)
	countByTask := make(map[uint]int)

	for _, entry := range entries {
		totalsByTask[entry.TaskID] += entry.Duration
		countByTask[entry.TaskID]++
	}

	results := make([]TaskProductivityEntry, 0, len(tasks))

	for _, task := range tasks {
		seconds := totalsByTask[task.ID]

		results = append(results, TaskProductivityEntry{
			TaskID:       task.ID,
			Title:        task.Title,
			Status:       task.Status,
			Priority:     task.Priority,
			TotalSeconds: seconds,
			TotalHours:   roundToTwoDecimals(float64(seconds) / 3600),
			EntryCount:   countByTask[task.ID],
		})
	}

	sortTaskProductivityDesc(results)

	return results, nil
}

// =========================
// USER ACTIVITY
// =========================

type UserActivity struct {
	TotalProjects    int64   `json:"total_projects"`
	TotalTasks       int64   `json:"total_tasks"`
	CompletedTasks   int64   `json:"completed_tasks"`
	TotalTimeSeconds int64   `json:"total_time_seconds"`
	TotalTimeHours   float64 `json:"total_time_hours"`
	HasActiveTimer   bool    `json:"has_active_timer"`
}

// UserActivity returns a high-level activity summary for the
// authenticated user across all of their projects.
func (s *ReportService) UserActivity(
	userID uint,
) (*UserActivity, error) {

	var totalProjects int64
	if err := s.DB.Model(&models.Project{}).
		Where("user_id = ?", userID).
		Count(&totalProjects).Error; err != nil {
		return nil, err
	}

	var totalTasks int64
	if err := s.DB.Model(&models.Task{}).
		Where("user_id = ?", userID).
		Count(&totalTasks).Error; err != nil {
		return nil, err
	}

	var completedTasks int64
	if err := s.DB.Model(&models.Task{}).
		Where("user_id = ? AND status = ?", userID, "completed").
		Count(&completedTasks).Error; err != nil {
		return nil, err
	}

	var entries []models.TimeEntry
	if err := s.DB.
		Where("user_id = ?", userID).
		Find(&entries).Error; err != nil {
		return nil, err
	}

	var totalSeconds int64
	hasActive := false

	for _, entry := range entries {
		totalSeconds += entry.Duration

		if entry.EndTime == nil {
			hasActive = true
		}
	}

	return &UserActivity{
		TotalProjects:    totalProjects,
		TotalTasks:       totalTasks,
		CompletedTasks:   completedTasks,
		TotalTimeSeconds: totalSeconds,
		TotalTimeHours:   roundToTwoDecimals(float64(totalSeconds) / 3600),
		HasActiveTimer:   hasActive,
	}, nil
}

// =========================
// DASHBOARD
// =========================

type DashboardReport struct {
	Activity      *UserActivity      `json:"activity"`
	RecentEntries []models.TimeEntry `json:"recent_time_entries"`
	TasksByStatus map[string]int64   `json:"tasks_by_status"`
}

// Dashboard combines activity totals, a task-status breakdown, and the
// most recent time entries into a single overview payload.
func (s *ReportService) Dashboard(
	userID uint,
) (*DashboardReport, error) {

	activity, err := s.UserActivity(userID)

	if err != nil {
		return nil, err
	}

	statuses := []string{"todo", "in_progress", "completed"}
	tasksByStatus := make(map[string]int64)

	for _, status := range statuses {
		var count int64

		if err := s.DB.Model(&models.Task{}).
			Where("user_id = ? AND status = ?", userID, status).
			Count(&count).Error; err != nil {
			return nil, err
		}

		tasksByStatus[status] = count
	}

	var recentEntries []models.TimeEntry

	if err := s.DB.
		Where("user_id = ?", userID).
		Order("start_time DESC").
		Limit(5).
		Find(&recentEntries).Error; err != nil {
		return nil, err
	}

	return &DashboardReport{
		Activity:      activity,
		RecentEntries: recentEntries,
		TasksByStatus: tasksByStatus,
	}, nil
}

// =========================
// HELPERS
// =========================

func roundToTwoDecimals(value float64) float64 {
	return float64(int64(value*100+0.5)) / 100
}

// sortTaskProductivityDesc sorts by total time spent, descending.
// Implemented manually (insertion sort) since the slice is small and
// bounded by the user's task count, avoiding an extra import.
func sortTaskProductivityDesc(entries []TaskProductivityEntry) {
	for i := 1; i < len(entries); i++ {
		j := i

		for j > 0 && entries[j-1].TotalSeconds < entries[j].TotalSeconds {
			entries[j-1], entries[j] = entries[j], entries[j-1]
			j--
		}
	}
}
