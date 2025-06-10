package handlers

import (
	"fmt"
	"net/smtp"
	"os"
	"regexp"
	"strings"
)

// Send sends a minimal HTML email using <subject> and <body> blocks in the template.
func Email(toEmail string, templatePath string, replacements map[string]string) error {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read email template: %w", err)
	}

	raw := string(content)

	// Extract <subject>...</subject>
	subjectMatch := regexp.MustCompile(`(?s)<subject>(.*?)</subject>`).FindStringSubmatch(raw)
	if len(subjectMatch) < 2 {
		return fmt.Errorf("missing <subject> block in template")
	}
	subject := strings.TrimSpace(subjectMatch[1])

	// Extract <body>...</body>
	bodyMatch := regexp.MustCompile(`(?s)<body>(.*?)</body>`).FindStringSubmatch(raw)
	if len(bodyMatch) < 2 {
		return fmt.Errorf("missing <body> block in template")
	}
	bodyHTML := strings.TrimSpace(bodyMatch[1])

	// Enforce: No inline styles (to preserve native email client rendering)
	if strings.Contains(bodyHTML, "style=") {
		return fmt.Errorf("inline styling is not allowed to preserve native rendering")
	}

	// Apply replacements
	for key, val := range replacements {
		subject = strings.ReplaceAll(subject, key, val)
		bodyHTML = strings.ReplaceAll(bodyHTML, key, val)
	}

	// Compose raw HTML email — no multipart, no styling
	msg := "To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		bodyHTML

	// SMTP setup
	auth := smtp.PlainAuth("", cfg.SMTPUser, cfg.SMTPPass, cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)

	// Send email
	if err := smtp.SendMail(addr, auth, cfg.SMTPFrom, []string{toEmail}, []byte(msg)); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}
