package service

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) error {
	from := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMPT_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	// Create a new message
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Create a new dialer with timeout settings
	d := gomail.NewDialer(smtpHost, smtpPort, from, password)
	d.SSL = false
	d.TLSConfig = nil

	// Try to send the email with retries
	maxRetries := 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if err := d.DialAndSend(m); err != nil {
			lastErr = err
			// Wait before retrying
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return nil
	}

	return fmt.Errorf("failed to send email after %d attempts: %v", maxRetries, lastErr)
}

func SendRecommendationJob(to string, subject string, body string) error {
	// Send email directly using SendEmail
	err := SendEmail(to, subject, body)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}
