// Package email provides email service infrastructure
package smtp

import (
	"github.com/samber/lo"
	"gopkg.in/gomail.v2"
)

// Client represents the email service client
type Client struct {
	from     string
	fromName string
	dialer   *gomail.Dialer
}

// Message represents an email message
type Message struct {
	// Recipient email addresses
	To []string
	// CC email addresses
	CC []string
	// BCC email addresses
	BCC []string
	// Email subject
	Subject string
	// Email body (plain text)
	Body string
	// Email body (HTML)
	HTMLBody string
	// Attachments (file paths)
	Attachments []string
}

// SendMail sends an email message
func (c *Client) SendMail(msg *Message) error {
	m := gomail.NewMessage()

	// Set sender
	if c.fromName != "" {
		m.SetHeader("From", m.FormatAddress(c.from, c.fromName))
	} else {
		m.SetHeader("From", c.from)
	}

	// Set recipients
	m.SetHeader("To", msg.To...)

	// Set CC if provided
	if len(msg.CC) > 0 {
		m.SetHeader("Cc", msg.CC...)
	}

	// Set BCC if provided
	if len(msg.BCC) > 0 {
		m.SetHeader("Bcc", msg.BCC...)
	}

	// Set subject
	m.SetHeader("Subject", msg.Subject)

	// Set body (prefer HTML if both are provided)
	if msg.HTMLBody != "" {
		m.SetBody("text/html", msg.HTMLBody)
		if msg.Body != "" {
			m.AddAlternative("text/plain", msg.Body)
		}
	} else {
		m.SetBody("text/plain", msg.Body)
	}

	// Attach files if provided
	for _, attachment := range msg.Attachments {
		m.Attach(attachment)
	}

	// Send the email
	if err := c.dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

// SendTextMail sends a simple text email to a single recipient
func (c *Client) SendTextMail(to, subject, body string) error {
	return c.SendMail(&Message{
		To:      []string{to},
		Subject: subject,
		Body:    body,
	})
}

// SendHTMLMail sends an HTML email to a single recipient
func (c *Client) SendHTMLMail(to, subject, htmlBody string) error {
	return c.SendMail(&Message{
		To:       []string{to},
		Subject:  subject,
		HTMLBody: htmlBody,
	})
}

// SendBatchMail sends the same email to multiple recipients
func (c *Client) SendBatchMail(to []string, subject, body string) error {
	return c.SendMail(&Message{
		To:      to,
		Subject: subject,
		Body:    body,
	})
}

// SendBatchHTMLMail sends an HTML email to multiple recipients
func (c *Client) SendBatchHTMLMail(to []string, subject, htmlBody string) error {
	return c.SendMail(&Message{
		To:       to,
		Subject:  subject,
		HTMLBody: htmlBody,
	})
}

// Ping tests the SMTP connection
func (c *Client) Ping() {
	sender := lo.Must1(c.dialer.Dial())
	defer sender.Close()
}
