package handlers

import (
	"bytes"
	"context"
	"crowdfund/backend/models"
	"crowdfund/backend/services"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Mock UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetUserByUsername(username string) (models.User, error) {
	args := m.Called(username)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserService) GetUserByID(id uint) (models.User, error) {
	args := m.Called(id)
	return args.Get(0).(models.User), args.Error(1)
}

// Mock CacheService
type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) Get(ctx context.Context, key string, value interface{}) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockCacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockCacheService) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateProjectCache(projectID uint64) {
	m.Called(projectID)
}

func (m *MockCacheService) InvalidateUserCache(userID uint) {
	m.Called(userID)
}

// TestLogin_Success tests successful login
func TestLogin_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	// Create mock services
	mockUserService := new(MockUserService)
	mockCacheService := new(MockCacheService)
	
	// Create a mock user with hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockUser := models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}
	
	// Set up expectations
	mockUserService.On("GetUserByUsername", "testuser").Return(mockUser, nil)
	
	// Create handler with mock services
	// Convert mockUserService to the expected type
	userServicePtr := (*services.UserService)(nil)
	handler := NewUserHandlers(userServicePtr, mockCacheService)
	
	// Create a test router
	router := gin.New()
	router.POST("/users/login", handler.Login)
	
	// Create a login request
	loginData := map[string]string{
		"username": "testuser",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(loginData)
	
	// Create a test request
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Perform the request
	router.ServeHTTP(w, req)
	
	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Parse the response
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	
	// Check that we got a token
	assert.NotNil(t, response["token"])
	
	// Check that we got the user
	user, ok := response["user"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, float64(1), user["id"])
	assert.Equal(t, "testuser", user["username"])
	assert.Equal(t, "test@example.com", user["email"])
	
	// Verify expectations
	mockUserService.AssertExpectations(t)
}

// TestLogin_InvalidCredentials tests login with invalid credentials
func TestLogin_InvalidCredentials(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	// Create mock services
	mockUserService := new(MockUserService)
	mockCacheService := new(MockCacheService)
	
	// Set up expectations
	mockUserService.On("GetUserByUsername", "nonexistentuser").Return(models.User{}, errors.New("user not found"))
	
	// Create handler with mock services
	// Convert mockUserService to the expected type
	userServicePtr := (*services.UserService)(nil)
	handler := NewUserHandlers(userServicePtr, mockCacheService)
	
	// Create a test router
	router := gin.New()
	router.POST("/users/login", handler.Login)
	
	// Create a login request
	loginData := map[string]string{
		"username": "nonexistentuser",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(loginData)
	
	// Create a test request
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Perform the request
	router.ServeHTTP(w, req)
	
	// Check the response
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	// Verify expectations
	mockUserService.AssertExpectations(t)
}

// TestLogin_WrongPassword tests login with wrong password
func TestLogin_WrongPassword(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	// Create mock services
	mockUserService := new(MockUserService)
	mockCacheService := new(MockCacheService)
	
	// Create a mock user with hashed password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	mockUser := models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}
	
	// Set up expectations
	mockUserService.On("GetUserByUsername", "testuser").Return(mockUser, nil)
	
	// Create handler with mock services
	// Convert mockUserService to the expected type
	userServicePtr := (*services.UserService)(nil)
	handler := NewUserHandlers(userServicePtr, mockCacheService)
	
	// Create a test router
	router := gin.New()
	router.POST("/users/login", handler.Login)
	
	// Create a login request with wrong password
	loginData := map[string]string{
		"username": "testuser",
		"password": "wrongpassword",
	}
	jsonData, _ := json.Marshal(loginData)
	
	// Create a test request
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Perform the request
	router.ServeHTTP(w, req)
	
	// Check the response
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	// Verify expectations
	mockUserService.AssertExpectations(t)
}

// TestLogin_InvalidRequestFormat tests login with invalid request format
func TestLogin_InvalidRequestFormat(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	
	// Create mock services
	mockUserService := new(MockUserService)
	mockCacheService := new(MockCacheService)
	
	// Create handler with mock services
	// Convert mockUserService to the expected type
	userServicePtr := (*services.UserService)(nil)
	handler := NewUserHandlers(userServicePtr, mockCacheService)
	
	// Create a test router
	router := gin.New()
	router.POST("/users/login", handler.Login)
	
	// Create an invalid JSON request
	invalidJSON := []byte(`{"username": "testuser", "password":}`)
	
	// Create a test request
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	
	// Create a response recorder
	w := httptest.NewRecorder()
	
	// Perform the request
	router.ServeHTTP(w, req)
	
	// Check the response
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
