/*
 * Email stuff to send reports by mail.
 */

package main

import (
	"fmt"
	"net/smtp"
	"time"
)

// timeFormat is the format for the date-time that appears in mails.
const timeFormat = "2006-01-02 15:04"

// EmailSubject builds the email subject from the build report.
func EmailSubject(build Build, start time.Time) string {
	return fmt.Sprintf("%s build on %s was a %s", build.Module.Name, start.Format(timeFormat), build)
}

// EmailMessage builds the body of the email.
func EmailMessage(build Build, start time.Time, duration time.Duration,
	email EmailConfig, subject string) string {
	message := fmt.Sprintf("From: %s\n", email.Sender)
	message += fmt.Sprintf("To: %s\n", email.Recipient)
	message += fmt.Sprintf("Subject: %s\n\n", subject)
	message += fmt.Sprintf("%s build on %s was a %s.\n", build.Module.Name, start.Format(timeFormat), build.String())
	message += fmt.Sprintf("Done in %s", duration)
	if !build.Success {
		message += fmt.Sprintf("\n\n===================================\n")
		message += fmt.Sprintf("ERROR:")
		message += fmt.Sprintf("\n-----------------------------------\n")
		message += fmt.Sprintf(build.Output)
		message += fmt.Sprintf("\n-----------------------------------\n")
	}
	message += "\n--\ncontinuum"
	return message
}

// SendEmail sends the report email after the build.
func SendEmail(build Build, start time.Time, duration time.Duration, email EmailConfig) {
	if email.SmtpHost == "" {
		return
	}
	if build.SendEmail(email) {
		fmt.Print("Sending email... ")
		subject := EmailSubject(build, start)
		message := EmailMessage(build, start, duration, email, subject)
		err := smtp.SendMail(email.SmtpHost, nil, email.Sender,
			[]string{email.Recipient}, []byte(message))
		if err != nil {
			fmt.Println("ERROR")
		} else {
			fmt.Println("OK")
		}
	}
}
