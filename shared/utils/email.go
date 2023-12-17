package utils

import (
	"crypto/tls"
	"fmt"
	"github.com/i-Things/things/shared/conf"
	"github.com/jordan-wright/email"
	"net/smtp"
)

func SenEmail(c conf.Email, to []string, subject string, body string) error {
	if c.Port == 0 {
		c.Port = 465
	}
	auth := smtp.PlainAuth("", c.From, c.Secret, c.Host)
	e := email.NewEmail()
	if c.Nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", c.Nickname, c.From)
	} else {
		e.From = c.From
	}
	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)
	var err error
	hostAddr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	if c.IsSSL {
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: c.Host})
	} else {
		err = e.Send(hostAddr, auth)
	}
	return err
}
