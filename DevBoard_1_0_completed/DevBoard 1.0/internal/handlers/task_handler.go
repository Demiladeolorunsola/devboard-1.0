package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Demiladeolorunsola/devboard-1.0/internal/services"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	Service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{
		Service: service,
	}
}

// =========================
// CREATE TASK
// =========================

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

func (h *TaskHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	projectID, err := parseTaskID(c.Param("project_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	var request CreateTaskRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	task, err := h.Service.Create(
		userID,
		projectID,
		request.Title,
		request.Description,
		request.Priority,
	)

	if err != nil {
		if err.Error() == "project not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "task created successfully",
		"task":    task,
	})
}

// =========================
// GET PROJECT TASKS
// =========================

func (h *TaskHandler) GetByProject(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	projectID, err := parseTaskID(c.Param("project_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	tasks, err := h.Service.GetByProject(
		userID,
		projectID,
	)

	if err != nil {
		if err.Error() == "project not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve tasks",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})
}

// =========================
// GET SINGLE TASK
// =========================

func (h *TaskHandler) GetByID(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	taskID, err := parseTaskID(c.Param("task_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}

	task, err := h.Service.GetByID(
		userID,
		taskID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task": task,
	})
}

// =========================
// UPDATE TASK
// =========================

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
}

func (h *TaskHandler) Update(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	taskID, err := parseTaskID(c.Param("task_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}

	var request UpdateTaskRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	task, err := h.Service.Update(
		userID,
		taskID,
		request.Title,
		request.Description,
		request.Status,
		request.Priority,
	)

	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "task updated successfully",
		"task":    task,
	})
}

// =========================
// DELETE TASK
// =========================

func (h *TaskHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	taskID, err := parseTaskID(c.Param("task_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}

	if err := h.Service.Delete(userID, taskID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "task deleted successfully",
	})
}

// =========================
// PARSE TASK ID
// =========================

func parseTaskID(value string) (uint, error) {
	id, err := strconv.ParseUint(value, 10, 64)

	if err != nil || id == 0 {
		return 0, fmt.Errorf("invalid id")
	}

	return uint(id), nil
}
