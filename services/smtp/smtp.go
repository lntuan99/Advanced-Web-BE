package smtp

import (
	"advanced-web.hcmus/config"
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

var (
	auth smtp.Auth
)

//Request struct
type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

func Initialize() {
	auth = smtp.PlainAuth(
		"",
		config.Config.SMTPUsername,
		config.Config.SMTPPassword,
		config.Config.SMTPHost)

	fmt.Println("Connected to Mail Service!")
}

func NewRequest(to []string, subject, body string) *Request {
	return &Request{
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (r *Request) SendEmail() (bool, error) {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: " + r.subject + "!\n"
	msg := []byte(subject + mime + "\n" + r.body)
	addr := fmt.Sprintf("smtp.gmail.com:%v", config.Config.SMTPPort)

	if err := smtp.SendMail(addr, auth, config.Config.SMTPUsername, r.to, msg); err != nil {
		return false, err
	}
	return true, nil
}

func (r *Request) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}