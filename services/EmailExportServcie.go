package services

import (
	"DB_race/models"
	"errors"
	"fmt"
	"net/smtp"
)

type EmailExportService struct {
	host, port, username, password string
}

func (service *EmailExportService) SetCredentials(host, port, username, password string) {
	service.host = host
	service.port = port
	service.username = username
	service.password = password
}

func (service *EmailExportService) Export(article models.Article) (path string, err error) {
	if service.host == "" {
		return "", errors.New("email credentials aren't set")
	}

	sender := "gosender@test.mail"
	receiver := []string{"goreceiver@test.mail"}

	subject := fmt.Sprintf("Exported article %d", article.Id)
	data := fmt.Sprintf("Title: %s\nUrl: %s\nCreated at: %s", article.Title, article.URL, article.CreatedAt)
	byteData := []byte("To: goreceiver@test.mail" +
		"\n\r" +
		"Subject: " + subject +
		"\n\r" +
		data + "\r\n")

	auth := smtp.PlainAuth("", service.username, service.password, service.host)
	smtpErr := smtp.SendMail(service.host+":"+service.port, auth, sender, receiver, byteData)

	if smtpErr != nil {
		return "", smtpErr
	}

	return subject, nil
}
