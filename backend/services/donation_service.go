package services

import (
	"crowdfund/backend/models"
	"log"
	"sync"

	"gorm.io/gorm"
)

type DonationService struct {
	db            *gorm.DB
	emailService  *EmailService
	donationTasks chan<- models.Donation
}

func NewDonationService(db *gorm.DB, emailService *EmailService, donationTasks chan<- models.Donation) *DonationService {
	return &DonationService{db: db, emailService: emailService, donationTasks: donationTasks}
}

func (s *DonationService) CreateDonation(donation models.Donation) error {
	s.donationTasks <- donation // Send to worker pool
	return nil
}

func DonationWorker(id int, tasks <-chan models.Donation, wg *sync.WaitGroup, db *gorm.DB, emailService *EmailService) {
	defer wg.Done()
	for task := range tasks {
		log.Printf("Worker %d processing donation: %v", id, task)
		if err := db.Create(&task).Error; err != nil {
			log.Printf("Error saving donation: %v", err)
			continue
		}
		// Send email confirmation
		emailService.SendDonationConfirmation(task)
	}
}

func (s *DonationService) GetDonationsByProjectID(projectID uint64) ([]models.Donation, error) {
	var donations []models.Donation
	err := s.db.Where("project_id = ?", projectID).Find(&donations).Error
	return donations, err
}
