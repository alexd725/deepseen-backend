package utilities

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	. "deepseen-backend/database/schemas"
)

// Send an email via SendGrid
func SendEmail(user *User, subject string, formattedTemplate Template) {
	from := mail.NewEmail("Deepseen", os.Getenv("SENDGRID_FROM_ADDRESS"))
	to := mail.NewEmail(user.FirstName+" "+user.LastName, user.Email)

	message := mail.NewSingleEmail(
		from,
		subject,
		to,
		formattedTemplate.Plain,
		formattedTemplate.Html,
	)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, emailError := client.Send(message)

	// show log
	if emailError == nil {
		fmt.Println("-- SENDGRID: sent to", user.Email, "[", response.StatusCode, response.Body, "]")
	} else {
		fmt.Println("-- SENDGRID: error", emailError, "[", user.Email, "]")
	}
}
