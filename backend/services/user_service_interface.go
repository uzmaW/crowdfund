package services

import (
	"crowdfund/backend/models"
)

// UserServiceInterface defines the interface for user service operations
type UserServiceInterface interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (models.User, error)
	GetUserByID(id uint) (models.User, error)
}

// Ensure UserService implements UserServiceInterface
var _ UserServiceInterface = (*UserService)(nil)
