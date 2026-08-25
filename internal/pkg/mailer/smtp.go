// internal/pkg/mailer/smtp.go
package mailer

import (
	"context"
	"fmt"
	"time"

	"github.com/DylanBergmann2502/go-maleficent/config"
	mail "github.com/wneessen/go-mail"
)

// Message contains the application-level fields needed to send an email.
type Message struct {
	To       []string
	Subject  string
	TextBody string
	HTMLBody string
}

// Mailer sends email through an SMTP server.
type Mailer struct {
	client   *mail.Client
	from     string
	fromName string
}

// New creates an SMTP mailer from the application mail configuration.
func New(config *config.MailConfig) (*Mailer, error) {
	options := []mail.Option{
		mail.WithPort(config.Port),
		mail.WithTimeout(15 * time.Second),
	}

	if config.TLS {
		options = append(options, mail.WithTLSPolicy(mail.TLSMandatory))
	} else {
		options = append(options, mail.WithTLSPolicy(mail.NoTLS))
	}

	if config.Username != "" || config.Password != "" {
		options = append(
			options,
			mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
			mail.WithUsername(config.Username),
			mail.WithPassword(config.Password),
		)
	}

	client, err := mail.NewClient(config.Host, options...)
	if err != nil {
		return nil, fmt.Errorf("create SMTP client: %w", err)
	}

	return &Mailer{
		client:   client,
		from:     config.From,
		fromName: config.FromName,
	}, nil
}

// Send delivers a message through the configured SMTP server.
func (m *Mailer) Send(ctx context.Context, message Message) error {
	if len(message.To) == 0 {
		return fmt.Errorf("email must have at least one recipient")
	}

	email := mail.NewMsg()
	var err error
	if m.fromName == "" {
		err = email.From(m.from)
	} else {
		err = email.FromFormat(m.fromName, m.from)
	}
	if err != nil {
		return fmt.Errorf("set email sender: %w", err)
	}
	if err := email.To(message.To...); err != nil {
		return fmt.Errorf("set email recipients: %w", err)
	}

	email.Subject(message.Subject)
	email.SetBodyString(mail.TypeTextPlain, message.TextBody)
	if message.HTMLBody != "" {
		email.AddAlternativeString(mail.TypeTextHTML, message.HTMLBody)
	}

	if err := m.client.DialAndSendWithContext(ctx, email); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}
