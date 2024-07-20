package gomail

import (
	"fmt"

	"github.com/wneessen/go-mail"
)

type MailClient struct {
	name   string
	email  string
	client *mail.Client
}

func NewMailClient(server string, username string, password string) (*MailClient, error) {
	c, err := mail.NewClient(server, mail.WithUsername(username), mail.WithPassword(password),
		mail.WithSMTPAuth(mail.SMTPAuthLogin), mail.WithPort(587))

	if err != nil {
		return nil, err
	}

	return &MailClient{
		client: c,
	}, nil
}

func (c *MailClient) SetFrom(name, email string) {
	c.name = name
	c.email = email
}

func (c *MailClient) Send(toName, toEmail, subject, msgBody string) error {
	m := mail.NewMsg()

	if err := m.EnvelopeFrom(fmt.Sprintf("noreply+%s@jobberapp.com", "jobber")); err != nil {
		return fmt.Errorf("failed to set ENVELOPE FROM address: %s", err)
	}

	if err := m.FromFormat(c.name, c.email); err != nil {
		return fmt.Errorf("failed to set formatted FROM address: %s", err)
	}

	if err := m.AddToFormat(toName, toEmail); err != nil {
		return fmt.Errorf("failed to set formatted TO address: %s", err)
	}

	m.SetMessageID()
	m.SetDate()
	m.Subject(subject)

	m.SetBodyString(mail.TypeTextHTML, msgBody)

	return c.client.DialAndSend(m)
}
