package mail

import (
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

type Mailer struct {
	Server   string
	User     string
	Password string
}

func (m *Mailer) SendEmail(subj, body string, receivers []string) error {
	addr := mail.Address{Name: "", Address: m.Server}
	host, _, _ := net.SplitHostPort(m.Server)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().String()))
	sb.WriteString(fmt.Sprintf("From: %s\r\n", m.User))
	sb.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(receivers, ",")))
	sb.WriteString(fmt.Sprintf("Bcc: %s\r\n", m.User))
	sb.WriteString(fmt.Sprintf("Subject: %s\r\n", subj))
	sb.WriteString("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n")
	sb.WriteString(body)

	msg := []byte(sb.String())
	auth := smtp.PlainAuth("", m.User, m.Password, host)

	receivers = append(receivers, m.User)
	return  smtp.SendMail(addr.Address, auth, m.User, receivers, msg)
}
