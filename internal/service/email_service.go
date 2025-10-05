package service

import (
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailService interface {
	SendWelcomeEmail(userEmail, userName string) error
	SendPasswordResetEmail(userEmail, userName, resetToken string) error
	SendEmail(toEmail, toName, subject, plainTextContent, htmlContent string) error
}

type emailService struct {
	client    *sendgrid.Client
	fromEmail string
	fromName  string
	templates *EmailTemplates
}

type EmailConfig struct {
	APIKey    string
	FromEmail string
	FromName  string
}

func NewEmailService(config EmailConfig) EmailService {
	client := sendgrid.NewSendClient(config.APIKey)
	return &emailService{
		client:    client,
		fromEmail: config.FromEmail,
		fromName:  config.FromName,
		templates: NewEmailTemplates(),
	}
}

// SendEmail sends a generic email with custom content
func (s *emailService) SendEmail(toEmail, toName, subject, plainTextContent, htmlContent string) error {
	from := mail.NewEmail(s.fromName, s.fromEmail)
	to := mail.NewEmail(toName, toEmail)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	response, err := s.client.Send(message)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode >= 400 {
		log.Printf("SendGrid error - Status: %d, Body: %s", response.StatusCode, response.Body)
		return fmt.Errorf("sendgrid returned error status: %d", response.StatusCode)
	}

	// Log success details (optional - you can remove this in production)
	fmt.Printf("Email sent successfully - Status: %d\n", response.StatusCode)

	return nil
}

// SendWelcomeEmail sends a welcome email using the template
func (s *emailService) SendWelcomeEmail(userEmail, userName string) error {
	subject, plainTextContent, htmlContent := s.templates.WelcomeEmailTemplate(userName)
	return s.SendEmail(userEmail, userName, subject, plainTextContent, htmlContent)
}

// SendPasswordResetEmail sends a password reset email using the template
func (s *emailService) SendPasswordResetEmail(userEmail, userName, resetToken string) error {
	subject, plainTextContent, htmlContent := s.templates.PasswordResetEmailTemplate(userName, resetToken)
	return s.SendEmail(userEmail, userName, subject, plainTextContent, htmlContent)
}
