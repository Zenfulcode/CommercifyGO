package service

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
)

// SMTPConfig contains SMTP server configuration
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// EmailService represents an email service implementation
type EmailService struct {
	config SMTPConfig
}

// NewEmailService creates a new instance of EmailService
func NewEmailService(config SMTPConfig) *EmailService {
	return &EmailService{
		config: config,
	}
}

// NewEmailServiceFromEnv creates a new instance of EmailService using environment variables
func NewEmailServiceFromEnv() *EmailService {
	return &EmailService{
		config: SMTPConfig{
			Host:     os.Getenv("SMTP_HOST"),
			Port:     os.Getenv("SMTP_PORT"),
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			From:     os.Getenv("SMTP_FROM"),
		},
	}
}

// SendEmail sends an email
func (s *EmailService) SendEmail(to, subject string, body []byte) error {
	// Prepare email headers
	var buf bytes.Buffer
	buf.WriteString("From: " + s.config.From + "\r\n")
	buf.WriteString("To: " + to + "\r\n")
	buf.WriteString("Subject: " + subject + "\r\n")
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
	buf.Write(body)

	// Connect to SMTP server
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)

	// Send email
	err := smtp.SendMail(
		addr,
		auth,
		s.config.From,
		[]string{to},
		buf.Bytes(),
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
