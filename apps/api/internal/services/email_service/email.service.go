package email_service

import (
	"context"

	"github.com/bored-bots/bored-template/internal/config"
	"github.com/wneessen/go-mail"
)

type Service struct {
	client *mail.Client
}

func New(smtpConfig *config.SMTPConfig) (*Service, error) {
	var opts []mail.Option
	opts = append(opts, mail.WithPort(smtpConfig.Port))
	opts = append(opts, mail.WithUsername(smtpConfig.Username))
	opts = append(opts, mail.WithPassword(smtpConfig.Password))
	opts = append(opts, mail.WithSMTPAuth(mail.SMTPAuthPlain))

	if smtpConfig.UseSSL {
		opts = append(opts, mail.WithSSL())
	} else {
		opts = append(opts, mail.WithTLSPolicy(mail.TLSOpportunistic))
	}

	client, err := mail.NewClient(smtpConfig.Host, opts...)
	if err != nil {
		return nil, err
	}

	return &Service{
		client: client,
	}, nil
}

type SendEmailInput struct {
	HTML     string
	Receiver string
	Sender   string
	Subject  string
}

func (s *Service) SendEmail(ctx context.Context, input SendEmailInput) error {
	m := mail.NewMsg()

	if err := m.From(input.Sender); err != nil {
		return err
	}

	if err := m.To(input.Receiver); err != nil {
		return err
	}

	m.Subject(input.Subject)
	m.SetBodyString(mail.TypeTextHTML, input.HTML)

	if err := s.client.DialAndSendWithContext(ctx, m); err != nil {
		return err
	}

	return nil
}
