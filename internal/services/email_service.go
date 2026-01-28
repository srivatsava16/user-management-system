package services

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendPasswordResetLink(email, token string) error {
	if email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_FROM")

	if username == "" || password == "" || smtpHost == "" || smtpPort == "" || from == "" {
		return fmt.Errorf("SMTP configuration is not set")
	}

	frontendURL := os.Getenv("FRONTEND_URL")

	if frontendURL == "" {
		return fmt.Errorf("FRONTEND_URL is not set")
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

	message := fmt.Sprintf(`Subject: Password Reset Request

	Hello,

	You requested a password reset. Click the link below to reset your password:

	%s

	This link will expire in 1 hour. 

	If you did not request this, please ignore this email.
	Best regards, ZXDS Platform `,
		resetLink)

	auth := smtp.PlainAuth("", username, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{email}, []byte(message))

	return err

}
