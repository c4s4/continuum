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
func EmailSubject(builds Builds, start time.Time) string {
	return fmt.Sprintf("Build on %s was a %s", start.Format(timeFormat), builds)
}

// EmailMessage builds the body of the email.
func EmailMessage(builds Builds, start time.Time, duration time.Duration,
	email EmailConfig, subject string) string {
	message := fmt.Sprintf("From: %s\n", email.Sender)
	message += fmt.Sprintf("To: %s\n", email.Recipient)
	message += fmt.Sprintf("Subject: %s\n\n", subject)
	message += fmt.Sprintf("Build on %s:\n\n", start.Format(timeFormat))
	for _, build := range builds {
		message += fmt.Sprintf("  %s: %s\n", build.Module.Name, build.String())
	}
	message += fmt.Sprintf("\nDone in %s\n", duration)
	message += builds.String()
	for _, build := range builds {
		if !build.Success {
			message += fmt.Sprintf("\n\n===================================\n")
			message += fmt.Sprintf(build.Module.Name)
			message += fmt.Sprintf("\n-----------------------------------\n")
			message += fmt.Sprintf(build.Output)
			message += fmt.Sprintf("\n-----------------------------------\n")
		}
	}
	message += "\n--\ngontinuum"
	return message
}

// SendEmail sends the report email after the build.
func SendEmail(builds Builds, start time.Time, duration time.Duration, email EmailConfig) {
	if !builds.Success() || (builds.Success() && email.Success) {
		subject := EmailSubject(builds, start)
		message := EmailMessage(builds, start, duration, email, subject)
		err := smtp.SendMail(email.SmtpHost, nil, email.Sender,
			[]string{email.Recipient}, []byte(message))
		if err != nil {
			panic(err)
		}
	}
}