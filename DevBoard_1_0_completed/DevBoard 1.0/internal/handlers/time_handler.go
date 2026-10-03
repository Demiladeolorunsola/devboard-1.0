package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Demiladeolorunsola/devboard-1.0/internal/services"
	"github.com/gin-gonic/gin"
)

type TimeHandler struct {
	Service *services.TimeService
}

func NewTimeHandler(service *services.TimeService) *TimeHandler {
	return &TimeHandler{
		Service: service,
	}
}

// =========================
// START TIMER
// =========================

type StartTimerRequest struct {
	ProjectID uint `json:"project_id"`
	TaskID    uint `json:"task_id"`
}

func (h *TimeHandler) Start(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var request StartTimerRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if request.ProjectID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "project_id is required",
		})
		return
	}

	if request.TaskID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "task_id is required",
		})
		return
	}

	timeEntry, err := h.Service.Start(
		userID,
		request.ProjectID,
		request.TaskID,
	)

	if err != nil {
		switch err.Error() {
		case "project not found", "task not found":
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "timer started successfully",
		"time_entry": timeEntry,
	})
}

// =========================
// STOP TIMER
// =========================

func (h *TimeHandler) Stop(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	timeEntryID, err := parseTimeID(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid time entry id",
		})
		return
	}

	timeEntry, err := h.Service.Stop(
		userID,
		timeEntryID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "timer stopped successfully",
		"time_entry": timeEntry,
	})
}

// =========================
// GET ALL TIME ENTRIES
// =========================

func (h *TimeHandler) GetAll(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	entries, err := h.Service.GetAll(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve time entries",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"time_entries": entries,
	})
}

// =========================
// GET PROJECT TIME
// =========================

func (h *TimeHandler) GetByProject(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	projectID, err := parseTimeID(c.Param("project_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	entries, err := h.Service.GetByProject(
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
			"error": "failed to retrieve project time",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"time_entries": entries,
	})
}

// =========================
// GET TASK TIME
// =========================

func (h *TimeHandler) GetByTask(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	taskID, err := parseTimeID(c.Param("task_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task id",
		})
		return
	}

	entries, err := h.Service.GetByTask(
		userID,
		taskID,
	)

	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve task time",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"time_entries": entries,
	})
}

// =========================
// PARSE TIME ID
// =========================

func parseTimeID(value string) (uint, error) {
	id, err := strconv.ParseUint(value, 10, 64)

	if err != nil || id == 0 {
		return 0, fmt.Errorf("invalid id")
	}

	return uint(id), nil
}
