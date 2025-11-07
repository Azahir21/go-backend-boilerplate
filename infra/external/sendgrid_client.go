package external

import (
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendGridConfig holds configuration for SendGrid email client.
type SendGridConfig struct {
	APIKey string `mapstructure:"api_key"`
	From   string `mapstructure:"from"`
}

// SendGridClient implements the EmailClient interface for SendGrid.
type SendGridClient struct {
	client *sendgrid.Client
	cfg    SendGridConfig
}

// NewSendGridClient creates a new SendGridClient instance.
func NewSendGridClient(cfg SendGridConfig) (EmailClient, error) {
	if cfg.APIKey == "" || cfg.From == "" {
		return nil, fmt.Errorf("SendGrid configuration (APIKey, From) cannot be empty")
	}
	return &SendGridClient{
		client: sendgrid.NewSendClient(cfg.APIKey),
		cfg:    cfg,
	}, nil
}

// SendEmail sends an email using SendGrid.
func (s *SendGridClient) SendEmail(to, subject, body string) error {
	from := mail.NewEmail("", s.cfg.From)
	toEmail := mail.NewEmail("", to)

	message := mail.NewSingleEmail(from, subject, toEmail, body, body)
	response, err := s.client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email via SendGrid: %w", err)
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("failed to send email via SendGrid, status code: %d, body: %s", response.StatusCode, response.Body)
	}

	return nil
}
