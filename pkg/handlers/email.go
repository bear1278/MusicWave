package handlers

import (
	"github.com/bear1278/MusicWave/configs"
	"gopkg.in/gomail.v2"
)

const (
	subject = "MusicWave Reset Password"
)

func SendMail(to, body string) error {
	m := gomail.NewMessage()
	cfg, err := configs.Init()
	if err != nil {
		return err
	}
	m.SetHeader("From", cfg.Email.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer("smtp.gmail.com", 587, cfg.Email.From, cfg.Email.Password) // Replace with your SMTP server details
	return d.DialAndSend(m)
}
