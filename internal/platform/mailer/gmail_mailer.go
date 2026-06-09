package mailer

import (
	"net/smtp"
	"os"
	"strings"
)

type Mailer interface {
	SendMail(to string, subject string, htmlBody string) error
}

type gmailMailer struct {
	fromAddress string
	appPassword string
}

func NewGmailMailer() Mailer {
	fromAddress := os.Getenv("GMAIL_ADDRESS")
	appPassword := os.Getenv("GMAIL_APP_PASSWORD")
	if fromAddress == "" || appPassword == "" {
		panic("mailer: GMAIL_ADDRESS and GMAIL_APP_PASSWORD must be set")
	}
	return &gmailMailer{
		fromAddress: fromAddress,
		appPassword: appPassword,
	}
}

func (m *gmailMailer) SendMail(to, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", m.fromAddress, m.appPassword, "smtp.gmail.com")

	var msg strings.Builder
	msg.WriteString("From: " + m.fromAddress + "\r\n")
	msg.WriteString("To: " + to + "\r\n")
	msg.WriteString("Subject: " + subject + "\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	return smtp.SendMail("smtp.gmail.com:587", auth, m.fromAddress, []string{to}, []byte(msg.String()))
}
