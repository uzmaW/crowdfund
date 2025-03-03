package handlers

import (
	"context"
	"crowdfund/backend/models"
	"crowdfund/backend/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ProjectHandlers struct {
	projectService *services.ProjectService
	cacheService   *services.CacheService
}

func NewProjectHandlers(projectService *services.ProjectService, cacheService *services.CacheService) *ProjectHandlers {
	return &ProjectHandlers{projectService: projectService, cacheService: cacheService}
}

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new project
// @Tags projects
// @Accept json
// @Produce json
// @Param project body models.Project true "Project details"
// @Security ApiKeyAuth
// @Success 201 {object} map[string]string{"message": "Project created successfully"}
// @Failure 400 {object} map[string]string{"error": "Invalid request"}
// @Failure 500 {object} map[string]string{"error": "Internal server error"}
// @Router /api/projects [post]
func (h *ProjectHandlers) CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userModel := user.(models.User)
	project.UserID = userModel.ID

	if err := h.projectService.CreateProject(&project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Project created successfully"})
}

// GetProject godoc
// @Summary Get a project by ID
// @Description Get a project by ID
// @Tags projects
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]string{"error": "Invalid project ID"}
// @Failure 500 {object} map[string]string{"error": "Internal server error"}
// @Router /api/projects/{id} [get]
func (h *ProjectHandlers) GetProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	ctx := context.Background()
	var project models.Project
	cacheKey := "project:" + strconv.FormatUint(id, 10)

	if err := h.cacheService.Get(ctx, cacheKey, &project); err == nil {
		c.JSON(http.StatusOK, project)
		return
	}

	project, err = h.projectService.GetProject(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.cacheService.Set(ctx, cacheKey, project, 1*time.Hour); err != nil {
		// Log error but don't fail request
	}

	c.JSON(http.StatusOK, project)
}

// UpdateProject godoc
// @Summary Update a project
// @Description Update a project
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param project body models.Project true "Updated project details"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string{"message": "Project updated successfully"}
// @Failure 400 {object} map[string]string{"error": "Invalid request"}
// @Failure 500 {object} map[string]string{"error": "Internal server error"}
// @Router /api/projects/{id} [put]
func (h *ProjectHandlers) UpdateProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.ID = uint(id)

	if err := h.projectService.UpdateProject(&project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.cacheService.InvalidateProjectCache(id)

	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully"})
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Delete a project
// @Tags projects
// @Param id path int true "Project ID"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string{"message": "Project deleted successfully"}
// @Failure 400 {object} map[string]string{"error": "Invalid project ID"}
// @Failure 500 {object} map[string]string{"error": "Internal server error"}
// @Router /api/projects/{id} [delete]
func (h *ProjectHandlers) DeleteProject(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	if err := h.projectService.DeleteProject(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.cacheService.InvalidateProjectCache(id)

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// ListProjects godoc
// @Summary List all projects
// @Description List all projects
// @Tags projects
// @Produce json
// @Success 200 {array} models.Project
// @Failure 500 {object} map[string]string{"error": "Internal server error"}
// @Router /api/projects [get]
func (h *ProjectHandlers) ListProjects(c *gin.Context) {
	projects, err := h.projectService.ListProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}
