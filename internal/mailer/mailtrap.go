package mailer

import (
	"bytes"
	"errors"
	"html/template"

	"gopkg.in/gomail.v2"
)

type mailtrapClient struct {
	fromEmail string
	apiKey    string
}

func NewMailtrapClient(apiKey, fromEmail string) (mailtrapClient, error) {
	if apiKey == "" {
		return mailtrapClient{}, errors.New("apiKey is required")
	}

	return mailtrapClient{
		apiKey:    apiKey,
		fromEmail: fromEmail,
	}, nil
}

func (m mailtrapClient) Send(templateFile, username, email string, data any, isSabdbox bool) (int, error) {
	// Template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject.tmpl", data)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body.tmpl", data)
	if err != nil {
		return -1, err
	}

	message := gomail.NewMessage()
	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())

	message.AddAlternative("text/html", body.String())

	dialer := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, "api", m.apiKey)

	if err := dialer.DialAndSend(message); err != nil {
		return -1, err
	}

	return 200, nil
}
