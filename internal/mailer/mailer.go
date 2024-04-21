package mailer

import (
	"github.com/go-mail/mail/v2"
	"time"
)

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func New(host string, port int, username, password, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

func (m Mailer) Send(recipient, subject, plainBody, htmlBody, token string) error {
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plainBody)
	msg.SetBody("token", token)
	msg.AddAlternative("text/html", htmlBody)

	err := m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil

}
