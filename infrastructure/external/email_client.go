package external

import (
	"fmt"

	"github.com/azahir21/go-backend-boilerplate/pkg/config"
	"github.com/sirupsen/logrus"
)

// EmailClient defines the interface for sending emails.
type EmailClient interface {
	SendEmail(to, subject, body string) error
}

// NewEmailClient creates a new email client based on configuration.
func NewEmailClient(log *logrus.Logger, cfg config.EmailConfig) (EmailClient, error) {
	switch cfg.Type {
	case "smtp":
		log.Info("Initializing SMTP Email Client...")
		smtpCfg := config.SmtpConfig{
			Host:     cfg.SMTP.Host,
			Port:     cfg.SMTP.Port,
			Username: cfg.SMTP.Username,
			Password: cfg.SMTP.Password,
			From:     cfg.SMTP.From,
		}
		return NewSmtpClient(smtpCfg)
	case "sendgrid":
		log.Info("Initializing SendGrid Email Client...")
		sendGridCfg := config.SendGridConfig{
			APIKey: cfg.SendGrid.APIKey,
			From:   cfg.SendGrid.From,
		}
		return NewSendGridClient(sendGridCfg)
	default:
		return nil, fmt.Errorf("unsupported email client type: %s", cfg.Type)
	}
}

// mockEmailClient is a placeholder for development.
type mockEmailClient struct{}

func (m *mockEmailClient) SendEmail(to, subject, body string) error {
	// fmt.Printf("Mock Email Sent to: %s, Subject: %s, Body: %s\n", to, subject, body)
	return nil
}