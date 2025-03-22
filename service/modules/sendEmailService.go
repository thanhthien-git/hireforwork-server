package service

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMPT_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(smtpHost, smtpPort, from, password)
	return d.DialAndSend(m)
}

func SendRecommendationJob(to string, subject string, body string) error {
	// Send email directly using SendEmail
	err := SendEmail(to, subject, body)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}
