package utilities

import (
	"os"
)

// Create a "Recovery" template
func CreateRecoveryTemplate(code, firstName, lastName string) Template {
	link := os.Getenv("BACKEND_URI") + "/api/auth/" + code
	return Template{
		Html: "<h1>Password recovery</h1>" +
			"<div>Hi, " + firstName + " " + lastName + "!</div>" +
			"<div>Here's your password recovery link:</div>" +
			"<div><a href='" + link + "'>" + link + "</a></div>",
		Plain: "Password recovery\nHi, " + firstName + " " + lastName + "!\n" +
			"Here's your password recovery link:\n" + link,
	}
}

// Create a "Welcome" template
func CreateWelcomeTemplate(firstName, lastName string) Template {
	return Template{
		Html: "<h1>Welcome to Deepseen!</h1>" +
			"<div>Hi, " + firstName + " " + lastName + "!</div>" +
			"<div>You can now use this email to sign in to your account in the desktop application.</div>",
		Plain: "Welcome to Deepseen!\nHi, " + firstName + " " + lastName + "!\n" +
			"You can sign in to your account in the desktop application.",
	}
}
