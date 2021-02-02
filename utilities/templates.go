package utilities

import (
	"fmt"
	"os"
)

// Reusable lines of text
var ignoreLine string = "You can safely ignore this message if you didn't create an account in Deepseen application."
var copyrightLine string = "Deepseen ©, all rights are reserved."

// Wrap HTML content
func wrapHtml(content string) string {
	return fmt.Sprintf(`
		<div style="background-color: black; color: white; padding: 48px;">
			%s
			<div style="font-size: 10px; margin-top: 36px;">
				<div>
					%s
				</div>
				<div>
					%s
				</div>
			</div>
		</div>
	`, content, ignoreLine, copyrightLine)
}

// Wrap plaintext content
func wrapPlain(content string) string {
	return fmt.Sprintf(`
		%s

		%s
		%s
	`, content, ignoreLine, copyrightLine)
}

// Create a "Recovery" template
func CreateRecoveryTemplate(code, firstName, lastName string) Template {
	line1 := "Account recovery"
	line2 := fmt.Sprintf("Hi, %s %s!", firstName, lastName)
	line3 := "Here's your account recovery link:"
	link := os.Getenv("FRONTEND_URI") + "/validate/" + code
	return Template{
		Html: wrapHtml(fmt.Sprintf(`
			<h1 style="color: turquoise; text-align: center;">
				%s
			</h1>
			<div style="font-size: 18px;">
				<div>
					%s
				</div>
				<div>
					%s
				</div>
				<div>
					<a href="%s" style="color: turquoise;">
						%s
					</a>
				</div>
			</div>
		`, line1, line2, line3, link, link)),
		Plain: wrapPlain(fmt.Sprintf(`
			%s
			
			%s
			%s
			%s
		`, line1, line2, line3, link)),
	}
}

// Create a "Welcome" template
func CreateWelcomeTemplate(firstName, lastName string) Template {
	line1 := "Welcome to Deepseen!"
	line2 := fmt.Sprintf("Hi, %s %s!", firstName, lastName)
	line3 := "You can now use this email address to sign in to your account in the desktop application."
	return Template{
		Html: wrapHtml(fmt.Sprintf(`
			<h1 style="color: turquoise; text-align: center;">
				%s
			</h1>
			<div style="font-size: 18px;">
				<div>
					%s
				</div>
				<div>
					%s
				</div>
			</div>
		`, line1, line2, line3)),
		Plain: wrapPlain(fmt.Sprintf(`
			%s
			%s
			%s
		`, line1, line2, line3)),
	}
}
