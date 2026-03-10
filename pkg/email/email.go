package email

import (
	"fmt"
	"net/smtp"
)

type EmailService struct {
	host     string
	port     string
	email    string
	password string
	appURL   string
}

func NewEmailService(host, port, email, password, appURL string) *EmailService {
	return &EmailService{
		host:     host,
		port:     port,
		email:    email,
		password: password,
		appURL:   appURL,
	}
}

func (e *EmailService) SendVerificationEmail(toEmail, token string) error {
	link := fmt.Sprintf("%s/auth/verify?token=%s", e.appURL, token)

	subject := "Verify your email"
	body := fmt.Sprintf(`
Hi!

Please verify your email by clicking the link below:

%s

This link will expire in 24 hours.

If you did not register, please ignore this email.
`, link)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		e.email, toEmail, subject, body)

	auth := smtp.PlainAuth("", e.email, e.password, e.host)
	addr := fmt.Sprintf("%s:%s", e.host, e.port)

	return smtp.SendMail(addr, auth, e.email, []string{toEmail}, []byte(msg))
}
