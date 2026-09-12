// Package mailer provides out-of-band delivery of secrets (e.g. password
// reset links) that must not be handed back to the requesting HTTP client.
package mailer

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

// sendTimeout bounds a single Send end-to-end so a slow or unreachable SMTP host can't hang the caller indefinitely.
const sendTimeout = 15 * time.Second

// Mailer sends a plain-text email; a returned error means only "could not deliver right now", nothing more specific.
type Mailer interface {
	Send(ctx context.Context, to, subject, body string) error
}

// FrontendURL returns the base URL of the frontend app, used to build links (e.g. a password-reset link) emailed to users.
func FrontendURL() string {
	if v := os.Getenv("FRONTEND_URL"); v != "" {
		return v
	}
	return "http://localhost:5173"
}

// NewFromEnv builds a Mailer from SMTP_* env vars, falling back to logging the message if SMTP_HOST is unset.
func NewFromEnv() Mailer {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		log.Println("mailer: SMTP_HOST not set — emails will be logged instead of sent")
		return &consoleMailer{}
	}

	port := os.Getenv("SMTP_PORT")
	if port == "" {
		port = "587"
	}

	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = "no-reply@openschool.local"
	}

	return &smtpMailer{
		host:     host,
		port:     port,
		username: os.Getenv("SMTP_USERNAME"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     from,
	}
}

// smtpMailer is the Mailer that actually delivers over SMTP once configured.
type smtpMailer struct {
	host, port, username, password, from string
}

// Send dials with a bounded deadline, uses implicit TLS on port 465 or opportunistic STARTTLS otherwise, and delivers one message.
func (m *smtpMailer) Send(ctx context.Context, to, subject, body string) error {
	deadline := time.Now().Add(sendTimeout)
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}

	addr := net.JoinHostPort(m.host, m.port)

	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("mailer: dial %s failed: %w", addr, err)
	}
	// One deadline on the connection bounds the whole conversation, since net/smtp has no context parameter of its own.
	if err := conn.SetDeadline(deadline); err != nil {
		conn.Close()
		return fmt.Errorf("mailer: setting deadline for %s failed: %w", addr, err)
	}

	if m.port == "465" {
		conn = tls.Client(conn, &tls.Config{ServerName: m.host, MinVersion: tls.VersionTLS12})
	}

	client, err := smtp.NewClient(conn, m.host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("mailer: SMTP handshake with %s failed: %w", addr, err)
	}
	defer client.Close()

	if m.port != "465" {
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(&tls.Config{ServerName: m.host, MinVersion: tls.VersionTLS12}); err != nil {
				return fmt.Errorf("mailer: STARTTLS with %s failed: %w", addr, err)
			}
		}
	}

	if m.username != "" {
		if err := client.Auth(smtp.PlainAuth("", m.username, m.password, m.host)); err != nil {
			return fmt.Errorf("mailer: authentication with %s failed: %w", addr, err)
		}
	}

	if err := client.Mail(m.from); err != nil {
		return fmt.Errorf("mailer: MAIL FROM rejected by %s: %w", addr, err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("mailer: RCPT TO rejected by %s: %w", addr, err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("mailer: DATA rejected by %s: %w", addr, err)
	}
	if _, err := w.Write([]byte(buildMessage(m.from, to, subject, body))); err != nil {
		return fmt.Errorf("mailer: writing message to %s failed: %w", addr, err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("mailer: finalizing message to %s failed: %w", addr, err)
	}

	return client.Quit()
}

// buildMessage assembles an RFC 5322 message with Date and Message-Id, since several mail servers spam-filter messages missing them.
func buildMessage(from, to, subject, body string) string {
	var idBytes [16]byte
	_, _ = rand.Read(idBytes[:])

	var b strings.Builder
	fmt.Fprintf(&b, "From: %s\r\n", from)
	fmt.Fprintf(&b, "To: %s\r\n", to)
	fmt.Fprintf(&b, "Subject: %s\r\n", subject)
	fmt.Fprintf(&b, "Date: %s\r\n", time.Now().UTC().Format(time.RFC1123Z))
	fmt.Fprintf(&b, "Message-Id: <%s@openschool>\r\n", hex.EncodeToString(idBytes[:]))
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	b.WriteString("\r\n")
	b.WriteString(body)
	b.WriteString("\r\n")
	return b.String()
}

// consoleMailer is the no-SMTP-configured fallback: it logs the email instead of delivering it.
type consoleMailer struct{}

// Send logs the message that would have been sent, since no SMTP server is configured.
func (consoleMailer) Send(_ context.Context, to, subject, body string) error {
	log.Printf("mailer: SMTP not configured, not sending email — to=%s subject=%q\n%s", to, subject, body)
	return nil
}
