package utilities

// Create a "Welcome" template
func CreateWelcomeTemplate(firstName, lastName string) Template {
	return Template{
		Html: "<h1>Welcome to Deepseen!</h1>" +
			"<div>Hi, " + firstName + " " + lastName + "!</div>" +
			"<div>You can sign in to your account in the desktop application.</div>",
		Plain: "Welcome to Deepseen!\nHi, " + firstName + " " + lastName + "!\n" +
			"You can sign in to your account in the desktop application.",
	}
}
