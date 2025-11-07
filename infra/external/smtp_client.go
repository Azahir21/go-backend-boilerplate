package external

import (
	"fmt"
	"net/smtp"

	"github.com/azahir21/go-backend-boilerplate/pkg/config"
)

// SmtpClient implements the EmailClient interface for SMTP.
type SmtpClient struct {
	cfg config.SmtpConfig
}

// NewSmtpClient creates a new SmtpClient instance.
func NewSmtpClient(cfg config.SmtpConfig) (EmailClient, error) {
	if cfg.Host == "" || cfg.Port == 0 || cfg.Username == "" || cfg.Password == "" || cfg.From == "" {
		return nil, fmt.Errorf("SMTP configuration (Host, Port, Username, Password, From) cannot be empty")
	}
	return &SmtpClient{cfg: cfg}, nil
}

// SendEmail sends an email using SMTP.
func (s *SmtpClient) SendEmail(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	err := smtp.SendMail(addr, auth, s.cfg.From, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email via SMTP: %w", err)
	}

	return nil
}

