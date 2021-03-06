package email

import (
	"fmt"
	"net/smtp"
	"originals/conf"
	"strings"

	"github.com/micro/go-log"
)

type Email struct {
	Receivers []string
	Subject   string
	Body      string
}

func SendMail(email *Email) error {
	log.Logf("Send email to %v\n", email.Receivers)
	smtpAddress := fmt.Sprintf(
		"%s:%d",
		conf.EmailConf.Host,
		conf.EmailConf.Port,
	)
	message := fmt.Sprintf(
		"To: %s\r\nFrom: %s\r\nSubject: %s\r\nContent-Type: %s\r\n\r\n%s",
		strings.Join(email.Receivers, ","),
		conf.EmailConf.Sender+"<"+conf.EmailConf.Sender+">",
		email.Subject,
		"text/html; charset=UTF-8",
		email.Body,
	)
	smtpAuth := smtp.PlainAuth(
		"",
		conf.EmailConf.UserName,
		conf.EmailConf.Password,
		conf.EmailConf.Host,
	)
	if err := smtp.SendMail(
		smtpAddress,
		smtpAuth,
		conf.EmailConf.UserName,
		email.Receivers,
		[]byte(message),
	); err != nil {
		log.Logf("Email send failed: %v\n", err)
		return err
	}
	return nil
}
