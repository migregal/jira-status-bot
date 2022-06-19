package mail

import (
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
)

type Mailer struct {
	Server   string
	User     string
	Password string
}

func (m *Mailer) SendEmail(subj, body string, receivers []string) error {
	addr := mail.Address{Name: "", Address: m.Server}
	host, _, _ := net.SplitHostPort(m.Server)
	msg := []byte(fmt.Sprintf("Subject: %s\n", subj) + body)
	auth := smtp.PlainAuth("", m.User, m.Password, host)

	receivers = append(receivers, m.User)
	return  smtp.SendMail(addr.Address, auth, m.User, receivers, msg)
}
