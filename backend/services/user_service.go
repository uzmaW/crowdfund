package services

import (
    "crowdfund/backend/models"

    "gorm.io/gorm"
)

type UserService struct {
    db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}

func (s *UserService) CreateUser(user *models.User) error {
    return s.db.Create(user).Error
}

func (s *UserService) GetUserByUsername(username string) (models.User, error) {
    var user models.User
    err := s.db.Where("username = ?", username).First(&user).Error
    return user, err
}

func (s *UserService) GetUserByID(id uint) (models.User, error) {
    var user models.User
    err := s.db.First(&user, id).Error
    return user, err
}