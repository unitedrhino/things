package main

import (
	"crypto/tls"
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
)

var (
	from     = "godlei6@qq.com"
	secret   = "hslpzkdqvzunbcid"
	host     = "smtp.qq.com"
	port     = 465
	nickname = "验证码机器人"
	to       = []string{"603785348@qq.com"}
	subject  = "邮箱登录验证"
	body     = "你的验证码是:123124"
	isSSL    = true
)

func main() {
	auth := smtp.PlainAuth("", from, secret, host)
	e := email.NewEmail()
	if nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", nickname, from)
	} else {
		e.From = from
	}
	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)
	var err error
	hostAddr := fmt.Sprintf("%s:%d", host, port)
	if isSSL {
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: host})
	} else {
		err = e.Send(hostAddr, auth)
	}
	fmt.Println(err)
}
