package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Demiladeolorunsola/devboard-1.0/internal/services"
	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	Service *services.ProjectService
}

func NewProjectHandler(service *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		Service: service,
	}
}

// =========================
// CREATE PROJECT
// =========================

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *ProjectHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var request CreateProjectRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	project, err := h.Service.Create(
		userID,
		request.Name,
		request.Description,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "project created successfully",
		"project": project,
	})
}

// =========================
// GET ALL PROJECTS
// =========================

func (h *ProjectHandler) GetAll(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	projects, err := h.Service.GetAll(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve projects",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"projects": projects,
	})
}

// =========================
// GET SINGLE PROJECT
// =========================

func (h *ProjectHandler) GetByID(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	projectID, err := parseID(c.Param("project_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	project, err := h.Service.GetByID(
		userID,
		projectID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"project": project,
	})
}

// =========================
// UPDATE PROJECT
// =========================

type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (h *ProjectHandler) Update(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	projectID, err := parseID(c.Param("project_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	var request UpdateProjectRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	project, err := h.Service.Update(
		userID,
		projectID,
		request.Name,
		request.Description,
		request.Status,
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

	c.JSON(http.StatusOK, gin.H{
		"message": "project updated successfully",
		"project": project,
	})
}

// =========================
// DELETE PROJECT
// =========================

func (h *ProjectHandler) Delete(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	projectID, err := parseID(c.Param("project_id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid project id",
		})
		return
	}

	if err := h.Service.Delete(userID, projectID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "project deleted successfully",
	})
}

// =========================
// PARSE PROJECT ID
// =========================

func parseID(value string) (uint, error) {
	id, err := strconv.ParseUint(value, 10, 64)

	if err != nil || id == 0 {
		return 0, fmt.Errorf("invalid id")
	}

	return uint(id), nil
}
