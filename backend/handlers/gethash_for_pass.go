package handlers

import (
	"crowdfund/backend/utils"
	"net/http"
	"github.com/gin-gonic/gin"
)

type PassHandlers struct {}

type PasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

func (h *PassHandlers) GetHashForPass(c *gin.Context) {
	var req PasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"hashedPassword": hashedPassword})
}
