package handlers

import (
        "crowdfund/backend/models"
        "crowdfund/backend/services"
        "net/http"
        "strconv"

        "github.com/gin-gonic/gin"
)

type DonationHandlers struct {
        donationService *services.DonationService
}

func NewDonationHandlers(donationService *services.DonationService) *DonationHandlers {
        return &DonationHandlers{donationService: donationService}
}

// CreateDonation godoc
// @Summary Create a new donation for a project
// @Description Create a new donation for a project
// @Tags donations
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param donation body models.Donation true "Donation details"
// @Security ApiKeyAuth
// @Success 201 {object} map[string]string{"message": "Donation created successfully"}
// @Failure 400 {object} map[string]string{"error": "Invalid request"}
// @Failure 500 {object} map[string]string{"error": "Internal server error"}
// @Router /api/projects/{id}/donations [post]
func (h *DonationHandlers) CreateDonation(c *gin.Context) {
        projectIDStr := c.Param("id")
        projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
                return
        }

        var donation models.Donation
        if err := c.ShouldBindJSON(&donation); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
        }

        user, exists := c.Get("user")
        if !exists {
                c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
                return
        }

        userModel := user.(models.User)

        donation.ProjectID = uint(projectID)
        donation.UserID = userModel.ID

        if err := h.donationService.CreateDonation(donation); err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusCreated, gin.H{"message": "Donation created successfully"})
}

// GetDonationsByProjectID godoc
// @Summary Get donations for a project
// @Description Get donations for a project
// @Tags donations
// @Produce json
// @Param id path int true "Project ID"
// @Security ApiKeyAuth
// @Success 200 {array} models.Donation
// @Failure 400 {object} map[string]string{"error": "Invalid request"}
// @Failure 500 {object} map[string]string{"error": "Internal server error"}
// @Router /api/projects/{id}/donations [get]
func (h *DonationHandlers) GetDonationsByProjectID(c *gin.Context) {
        projectIDStr := c.Param("id")
        projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
        if err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
                return
        }

        donations, err := h.donationService.GetDonationsByProjectID(projectID)
        if err != nil {
                c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
        }

        c.JSON(http.StatusOK, donations)
}