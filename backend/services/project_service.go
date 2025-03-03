package services

import (
    "crowdfund/backend/models"

    "gorm.io/gorm"
)

type ProjectService struct {
    db *gorm.DB
}

func NewProjectService(db *gorm.DB) *ProjectService {
    return &ProjectService{db: db}
}

func (s *ProjectService) CreateProject(project *models.Project) error {
    return s.db.Create(project).Error
}

func (s *ProjectService) GetProject(id uint64) (models.Project, error) {
    var project models.Project
    err := s.db.First(&project, id).Error
    return project, err
}

func (s *ProjectService) UpdateProject(project *models.Project) error {
    return s.db.Save(project).Error
}

func (s *ProjectService) DeleteProject(id uint64) error {
    return s.db.Delete(&models.Project{}, id).Error
}

func (s *ProjectService) ListProjects() ([]models.Project, error) {
    var projects []models.Project
    err := s.db.Find(&projects).Error
    return projects, err
}