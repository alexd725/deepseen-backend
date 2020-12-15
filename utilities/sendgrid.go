package utilities

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Send an email via SendGrid
func SendEmail(destination, firstName, lastName, subject, template string) {
	from := mail.NewEmail("Deepseen", os.Getenv("SENDGRID_FROM_ADDRESS"))
	to := mail.NewEmail(firstName+" "+lastName, destination)

	plainTextContent := "This is a test"
	htmlContent := "<strong>THIS is a test</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, emailError := client.Send(message)

	// show log
	if emailError == nil {
		fmt.Println("-- SENDGRID: sent to", destination, "[", response.StatusCode, response.Body, "]")
	} else {
		fmt.Println("-- SENDGRID: error", emailError, "[", destination, "]")
	}
}
