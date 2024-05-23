package email

import (
	"fmt"
	jemail "github.com/jordan-wright/email"
	"net/smtp"
)

type Config struct {
	Identity string `yaml:"identity"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	smtp     string
	port     int
}

type email struct {
	config *Config
	obj    *jemail.Email
}

func NewEmail(config *Config) *email {
	return &email{config: config, obj: jemail.NewEmail()}
}

func (e *email) QQmail() *email {
	e.config.smtp = "smtp.qq.com"
	e.config.port = 25
	return e
}

func (e *email) Gmail() *email {
	e.config.smtp = "smtp.gmail.com"
	e.config.port = 587
	return e
}

func (e *email) Send(parameter *EmailParameter) error {
	e.obj.From = parameter.From
	if e.obj.From == "" {
		e.obj.From = e.config.UserName
	}
	e.obj.To = parameter.To
	e.obj.Bcc = parameter.Bcc
	e.obj.Cc = parameter.Cc
	e.obj.HTML = parameter.HtmlBody
	e.obj.Text = parameter.TextBody
	e.obj.Subject = parameter.Subject
	return e.obj.Send(
		fmt.Sprintf("%s:%d", e.config.smtp, e.config.port),
		smtp.PlainAuth(e.config.Identity, e.config.UserName, e.config.Password, e.config.smtp),
	)
}
