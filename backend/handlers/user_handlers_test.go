package handlers

import (
	"bytes"
	"crowdfund/backend/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

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

func TestLogin_Success(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test-secret")

	// Create a test password hash
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	// Create a mock user
	mockUser := models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}

	// Create mock services
	mockUserService := new(MockUserService)
	mockCacheService := new(MockCacheService)

	// Set up expectations
	mockUserService.On("GetUserByUsername", "testuser").Return(mockUser, nil)

	// Create handler with mock services
	handler := NewUserHandlers(mockUserService, mockCacheService)

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

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// Assert response contains token and user
	assert.Contains(t, response, "token")
	assert.Contains(t, response, "user")

	// Verify expectations were met
	mockUserService.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create mock services
	mockUserService := new(MockUserService)
	mockCacheService := new(MockCacheService)

	// Set up expectations - user not found
	mockUserService.On("GetUserByUsername", "nonexistent").Return(models.User{}, errors.New("user not found"))

	// Create handler with mock services
	handler := NewUserHandlers(mockUserService, mockCacheService)

	// Create a test router
	router := gin.New()
	router.POST("/users/login", handler.Login)

	// Create a login request with invalid username
	loginData := map[string]string{
		"username": "nonexistent",
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

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Parse response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// Assert error message
	assert.Equal(t, "Invalid credentials", response["error"])

	// Verify expectations were met
	mockUserService.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create a test password hash
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)

	// Create a mock user
	mockUser := models.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}

	// Create mock services
	mockUserService := new(MockUserService)
	mockCacheService := new(MockCacheService)

	// Set up expectations
	mockUserService.On("GetUserByUsername", "testuser").Return(mockUser, nil)

	// Create handler with mock services
	handler := NewUserHandlers(mockUserService, mockCacheService)

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

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Parse response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// Assert error message
	assert.Equal(t, "Invalid credentials", response["error"])

	// Verify expectations were met
	mockUserService.AssertExpectations(t)
}

func TestLogin_InvalidRequestFormat(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)

	// Create mock services
	mockUserService := new(MockUserService)
	mockCacheService := new(MockCacheService)

	// Create handler with mock services
	handler := NewUserHandlers(mockUserService, mockCacheService)

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

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Verify no expectations were called
	mockUserService.AssertNotCalled(t, "GetUserByUsername")
}
