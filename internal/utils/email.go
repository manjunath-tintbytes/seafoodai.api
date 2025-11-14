package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailConfig struct {
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	From     string
}

func GetEmailConfig() *EmailConfig {
	return &EmailConfig{
		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPPort: os.Getenv("SMTP_PORT"),
		SMTPUser: os.Getenv("SMTP_USER"),
		SMTPPass: os.Getenv("SMTP_PASS"),
		From:     os.Getenv("SMTP_USER"),
	}
}

func SendPasswordResetEmail(to, token string) error {
	cfg := GetEmailConfig()

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

	subject := "Password Reset Request"
	body := fmt.Sprintf(`
Hello,

You requested to reset your password. Click the link below to reset your password:

%s

This link will expire in 1 hour.

If you didn't request this, please ignore this email.

Best regards,
Seafood AI Team
`, resetLink)

	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", cfg.From, to, subject, body)

	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)

	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)
	err := smtp.SendMail(addr, auth, cfg.From, []string{to}, []byte(message))
	println([]string{to})

	return err
}
