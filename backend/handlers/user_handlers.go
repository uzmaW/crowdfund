package handlers

import (
	"context"
	"crowdfund/backend/models"
	"crowdfund/backend/services"
	"crowdfund/backend/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserHandlers struct {
	userService  services.UserServiceInterface
	cacheService services.CacheServiceInterface
}

func NewUserHandlers(userService services.UserServiceInterface, cacheService services.CacheServiceInterface) *UserHandlers {
	return &UserHandlers{userService: userService, cacheService: cacheService}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User registration details"
// @Success 201 {object} map[string]string{"message": "User registered successfully"}
// @Failure 400 {object} map[string]string{"error": "Invalid request"}
// @Failure 500 {object} map[string]string{"error": "Internal server error"}
// @Router /users/register [post]
func (h *UserHandlers) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hashedPassword

	if err := h.userService.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login godoc
// @Summary Login a user
// @Description Login a user with username and password
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body models.LoginCredentials true "User login credentials"
// @Success 200 {object} map[string]interface{}{"token": "string", "user": models.User}
// @Failure 401 {object} map[string]string{"error": "Invalid credentials"}
// @Router /users/login [post]
func (h *UserHandlers) Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUserByUsername(loginData.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, os.Getenv("JWT_SECRET"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
}

// Profile godoc
// @Summary Get user profile
// @Description Get the profile information of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User "User profile information"
// @Failure 401 {object} map[string]string{"error": "Unauthorized"}
// @Failure 500 {object} map[string]string{"error": "Internal server error"}
// @Router /users/profile [get]
func (h *UserHandlers) Profile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userModel := user.(models.User)
	ctx := context.Background()
	cacheKey := "user:" + strconv.FormatUint(uint64(userModel.ID), 10)

	if err := h.cacheService.Get(ctx, cacheKey, &userModel); err == nil {
		c.JSON(http.StatusOK, userModel)
		return
	}

	userFromDb, err := h.userService.GetUserByID(userModel.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	if err := h.cacheService.Set(ctx, cacheKey, userFromDb, 1*time.Hour); err != nil {
		//log error
	}

	c.JSON(http.StatusOK, userFromDb)
}
