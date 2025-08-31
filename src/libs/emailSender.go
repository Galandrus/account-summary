package libs

import (
	"account-summary/src/config"
	"account-summary/src/models"
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
	"time"
)

type EmailSender interface {
	SendEmail(to string, subject string, summary models.TransactionSummary) error
}

type emailSender struct {
	from     string
	password string
	smtpHost string
	smtpPort string
}

func NewEmailSender(cfg *config.Config) EmailSender {
	return &emailSender{
		from:     cfg.EmailFrom,
		password: cfg.EmailPassword,
		smtpHost: cfg.SMTPHost,
		smtpPort: cfg.SMTPPort,
	}
}

type TemplateData struct {
	models.TransactionSummary
	GeneratedDate string
}

func (e *emailSender) SendEmail(to string, subject string, summary models.TransactionSummary) error {
	tmpl, err := template.ParseFiles("src/templates/summary.html")
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	templateData := TemplateData{
		TransactionSummary: summary,
		GeneratedDate:      time.Now().Format("January 2, 2006 at 3:04 PM"),
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, templateData)
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	message := e.createMIMEMessage(to, subject, body.String())

	auth := smtp.PlainAuth("", e.from, e.password, e.smtpHost)

	err = smtp.SendMail(
		fmt.Sprintf("%s:%s", e.smtpHost, e.smtpPort),
		auth,
		e.from,
		[]string{to},
		message,
	)
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	return nil
}

func (e *emailSender) createMIMEMessage(to, subject, htmlBody string) []byte {
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("Stori Card <%s>", e.from)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	var message strings.Builder
	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	message.WriteString("\r\n")
	message.WriteString(htmlBody)

	return []byte(message.String())
}
