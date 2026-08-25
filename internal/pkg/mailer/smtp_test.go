// internal/pkg/mailer/smtp_test.go
package mailer

import (
	"context"
	"testing"

	"github.com/DylanBergmann2502/go-maleficent/config"
	"github.com/stretchr/testify/require"
)

func TestNewBuildsSMTPMailer(t *testing.T) {
	mailer, err := New(&config.MailConfig{
		Host: "mailpit",
		Port: 1025,
		From: "noreply@example.local",
	})

	require.NoError(t, err)
	require.NotNil(t, mailer)
}

func TestSendRequiresRecipient(t *testing.T) {
	mailer := &Mailer{}

	err := mailer.Send(context.Background(), Message{})

	require.EqualError(t, err, "email must have at least one recipient")
}

func TestSendValidatesSenderBeforeConnecting(t *testing.T) {
	mailer := &Mailer{from: "invalid sender"}

	err := mailer.Send(context.Background(), Message{To: []string{"user@example.com"}})

	require.Error(t, err)
	require.Contains(t, err.Error(), "set email sender")
}
