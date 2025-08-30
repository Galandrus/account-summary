package libs

import (
	"account-summary/src/models"
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

type EmailSender interface {
	SendEmail(to string, subject string, summary models.TransactionSummary) error
	ServeHTML(summary models.TransactionSummary, port string) error
}

type emailSender struct {
}

func NewEmailSender() EmailSender {
	return &emailSender{}
}

func (e *emailSender) SendEmail(to string, subject string, summary models.TransactionSummary) error {
	// Parse the HTML template
	tmpl, err := template.ParseFiles("src/templates/email.html")
	if err != nil {
		return err
	}

	// Execute the template with the summary data
	var body bytes.Buffer
	err = tmpl.Execute(&body, summary)
	if err != nil {
		return err
	}

	// Here you would implement the actual email sending logic
	// For now, we'll just return nil as a placeholder

	return nil
}

// ServeHTML sirve el HTML parseado en un servidor web local
func (e *emailSender) ServeHTML(summary models.TransactionSummary, port string) error {
	if port == "" {
		port = "8080"
	}

	// Parse the HTML template
	tmpl, err := template.ParseFiles("src/templates/summary.html")
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	// Handler para servir el HTML
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := tmpl.Execute(w, summary)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		}
	})

	fmt.Printf("🌐 Servidor iniciado en http://localhost:%s\n", port)
	fmt.Println("Presiona Ctrl+C para detener el servidor...")

	return http.ListenAndServe(":"+port, nil)
}
