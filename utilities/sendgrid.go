package utilities

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Send an email via SendGrid
func SendEmail(address, template string) {
	from := mail.NewEmail("Deepseen", "test@example.com")
	subject := "Sending with Twilio SendGrid is Fun"
	to := mail.NewEmail("Example User", "test@example.com")
	plainTextContent := "and easy to do anywhere, even with Go"
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, emailError := client.Send(message)
	if emailError == nil {
		fmt.Println("-- mailer: sent to ", address, "[", response.StatusCode, response.Body, "]")
	}
}
