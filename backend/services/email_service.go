package services

import (
	"crowdfund/backend/models"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendDonationConfirmation(donation models.Donation) {
	e := email.NewEmail()
	e.From = "noreply@crowdfund.com"
	e.To = []string{"user@example.com"} // get user email from db.
	e.Subject = "Donation Confirmation"
	e.Text = []byte(fmt.Sprintf("Thank you for your donation of $%f", donation.Amount))

	auth := smtp.PlainAuth("", os.Getenv("MAILTRAP_USER"), os.Getenv("MAILTRAP_PASSWORD"), os.Getenv("MAILTRAP_HOST"))
	err := e.Send(os.Getenv("MAILTRAP_HOST")+":"+os.Getenv("MAILTRAP_PORT"), auth)
	if err != nil {
		log.Printf("Error sending email: %v", err)
	}
}
