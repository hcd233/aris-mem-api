// Package email provides email service infrastructure
package smtp

import (
	"crypto/tls"

	"github.com/hcd233/aris-mem-api/internal/config"
	"github.com/hcd233/aris-mem-api/internal/logger"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

var client *Client

// GetEmailClient returns the email client instance
//
//	return *Client
//	author centonhuang
//	update 2025-11-20 00:00:00
func GetEmailClient() *Client {
	return client
}

// InitSMTPClient initializes the email service
//
//	author centonhuang
//	update 2025-11-20 00:00:00
func InitSMTPClient() {
	dialer := gomail.NewDialer(config.SMTPHost, config.SMTPPort, config.SMTPUsername, config.SMTPPassword)

	// Configure TLS
	if config.SMTPTLS {
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         config.SMTPHost,
		}
	}

	client = &Client{
		from:     config.SMTPFrom,
		fromName: config.SMTPFromName,
		dialer:   dialer,
	}

	client.Ping()

	logger.Logger().Info("[SMTP] Connected to SMTP server", zap.String("host", config.SMTPHost), zap.Int("port", config.SMTPPort), zap.String("username", config.SMTPUsername), zap.Bool("tls", config.SMTPTLS), zap.String("from", config.SMTPFrom))
}
