package service

import (
	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) error {
	from := "hiresonarforwork@gmail.com"
	password := "poiy futi xurj mahu"
	smtpHost := "smtp.gmail.com"
	smtpPort := 587

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(smtpHost, smtpPort, from, password)
	return d.DialAndSend(m)
}
