package email

import (
	"fmt"
	"github.com/LMF709268224/titan-vps/node/config"
	"github.com/jordan-wright/email"
	"net/smtp"
	"net/textproto"
)

type Data struct {
	SendTo  string
	Subject string
	Tittle  string
	Content string
}

func SendEmail(cfg config.EmailConfig, data Data) error {
	message := &email.Email{
		To:      []string{data.SendTo},
		From:    fmt.Sprintf("%s", cfg.Name),
		Subject: data.Subject,
		Text:    []byte(data.Tittle),
		HTML:    []byte(data.Content),
		Headers: textproto.MIMEHeader{},
	}

	// smtp.PlainAuth：the first param can be empty，the second param should be the email account，the third param is the secret of the email
	addr := fmt.Sprintf("%s:%s", cfg.SMTPHost, cfg.SMTPPort)
	auth := smtp.PlainAuth("", cfg.Name, cfg.Password, cfg.SMTPHost)

	return message.Send(addr, auth)
}
