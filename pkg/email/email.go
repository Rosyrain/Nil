package email

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
	"time"
)

var (
	pool      *email.Pool
	maxClient int = 10
)

func SendEmail(from string, to []string, secret string, host string, nickname string, subject string, body string, port int, ssl bool) error {
	auth := smtp.PlainAuth("", from, secret, host)
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s<%s>", nickname, from)
	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)
	hostAddr := fmt.Sprintf("%s:%d", host, port)
	if ssl {
		return e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: host})
	}
	return e.Send(hostAddr, auth)
}

func SendEmailWithFile(from string, to []string, secret string, host string, nickname string, subject string, body string, port int, ssl bool, attach string) error {
	auth := smtp.PlainAuth("", from, secret, host)
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s<%s>", nickname, from)
	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)
	if attach != "" {
		_, _ = e.AttachFile(attach)
	}
	hostAddr := fmt.Sprintf("%s:%d", host, port)
	if ssl {
		return e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: host})
	}
	return e.Send(hostAddr, auth)
}

func SendEmailWithPool(to []string, from, secret, host, subject, body, nickname string, port int) (err error) {
	hostAddr := fmt.Sprintf("%s:%d", host, port)
	auth := smtp.PlainAuth("", from, secret, host)
	if pool == nil {
		pool, err = email.NewPool(hostAddr, maxClient, auth)
		if err != nil {
			log.Fatal(err)
		}
	}
	e := &email.Email{
		From:    fmt.Sprintf("%s<%s>", nickname, from),
		To:      to,
		Subject: subject,
		HTML:    []byte(body),
	}
	_ = tls.Config{InsecureSkipVerify: true}
	return pool.Send(e, 5*time.Second)
}

func SendEmailWithPoolAndFile(to []string, from, secret, host, subject, body, nickname string, port int, attach string) (err error) {
	hostAddr := fmt.Sprintf("%s:%d", host, port)
	auth := smtp.PlainAuth("", from, secret, host)
	if pool == nil {
		pool, err = email.NewPool(hostAddr, maxClient, auth)
		if err != nil {
			log.Fatal(err)
		}
	}
	e := &email.Email{
		From:    fmt.Sprintf("%s<%s>", nickname, from),
		To:      to,
		Subject: subject,
		Text:    []byte(body),
	}
	if attach != "" {
		_, _ = e.AttachFile(attach)
	}
	return pool.Send(e, 5*time.Second)
}
