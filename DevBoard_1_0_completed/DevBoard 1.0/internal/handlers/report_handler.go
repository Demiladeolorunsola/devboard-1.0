package handlers

import (
	"net/http"

	"github.com/Demiladeolorunsola/devboard-1.0/internal/services"
	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	Service *services.ReportService
}

func NewReportHandler(service *services.ReportService) *ReportHandler {
	return &ReportHandler{
		Service: service,
	}
}

// =========================
// PROJECT TIME SUMMARY
// =========================

func (h *ReportHandler) ProjectTimeSummary(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	projectID, err := parseID(c.Param("project_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	summary, err := h.Service.ProjectTimeSummary(userID, projectID)

	if err != nil {
		if err.Error() == "project not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate project time summary",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"report": summary,
	})
}

// =========================
// TASK PRODUCTIVITY
// =========================

func (h *ReportHandler) TaskProductivity(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	report, err := h.Service.TaskProductivity(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate task productivity report",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"report": report,
	})
}

// =========================
// USER ACTIVITY
// =========================

func (h *ReportHandler) UserActivity(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	activity, err := h.Service.UserActivity(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate activity report",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"report": activity,
	})
}

// =========================
// DASHBOARD
// =========================

func (h *ReportHandler) Dashboard(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	dashboard, err := h.Service.Dashboard(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate dashboard report",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"report": dashboard,
	})
}
