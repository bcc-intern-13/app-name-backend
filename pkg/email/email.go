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
	appURLFe string
}

func NewEmailService(host, port, email, password, appURL string, appURLFe string) *EmailService {
	return &EmailService{
		host:     host,
		port:     port,
		email:    email,
		password: password,
		appURL:   appURL,
		appURLFe: appURLFe,
	}
}

func (e *EmailService) SendVerificationEmail(toEmail, token string) error {
	link := fmt.Sprintf("%s/api/v1/auth/verify?token=%s", e.appURL, token)
	fmt.Println("APP URL:", e.appURL)

	fmt.Println(link)

	subject := "Verify Your Email Address"
	body := verificationEmailTemplate(toEmail, link)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		e.email, toEmail, subject, body)

	auth := smtp.PlainAuth("", e.email, e.password, e.host)
	addr := fmt.Sprintf("%s:%s", e.host, e.port)

	return smtp.SendMail(addr, auth, e.email, []string{toEmail}, []byte(msg))
}

func (e *EmailService) SendResetPasswordEmail(toEmail, token string) error {

	link := fmt.Sprintf("%s/reset-password?token=%s", e.appURLFe, token)
	fmt.Println("SENDING RESET LINK:", link)

	subject := "Reset Your WorkAble Password"
	body := resetPasswordEmailTemplate(toEmail, link)

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		e.email, toEmail, subject, body)

	auth := smtp.PlainAuth("", e.email, e.password, e.host)
	addr := fmt.Sprintf("%s:%s", e.host, e.port)

	return smtp.SendMail(addr, auth, e.email, []string{toEmail}, []byte(msg))
}
