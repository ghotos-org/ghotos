package mail

import (
	"crypto/tls"

	"io"
	"net/smtp"
	"os"
	"strconv"

	"github.com/jordan-wright/email"
)

type SMTPConfig struct {
	Host string
	Port int
	User string
	Pass string
}

type MailAttachment struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type SendMail struct {
	From        string
	To          []string
	Bcc         []string
	Cc          []string
	Subject     string
	Text        string
	Html        string
	ReadReceipt []string
	Attachments []MailAttachment
}

func GetConf(host string, port int, user string, pass string) SMTPConfig {
	return SMTPConfig{
		Host: host,
		Port: port,
		User: user,
		Pass: pass,
	}
}

// Send s1endet eine Email
func Send(sendMail SendMail, conf SMTPConfig) error {

	e := email.NewEmail()
	e.From = sendMail.From
	e.To = sendMail.To
	e.Subject = sendMail.Subject

	if sendMail.Bcc != nil {
		e.Bcc = sendMail.Bcc
	}
	if sendMail.Cc != nil {
		e.Cc = sendMail.Cc
	}
	if sendMail.ReadReceipt != nil {
		e.ReadReceipt = sendMail.ReadReceipt
	}
	if sendMail.Attachments != nil {

		sendMail.Text = sendMail.Text + "\n"

		if sendMail.Html != "" {
			sendMail.Html = sendMail.Html + "<br>"
		}

		for _, attachment := range sendMail.Attachments {

			osFile, err := os.Open(attachment.Path)
			if err != nil {
				return err
			}
			defer osFile.Close()

			var r io.Reader
			// this is doing a type conversion from *os.File to io.Reader
			r = (*os.File)(osFile)
			_, err = e.Attach(r, attachment.Name, attachment.Type)

			if err != nil {
				return err
			}
		}

	}

	e.Text = []byte(sendMail.Text)

	if sendMail.Html != "" {
		e.HTML = []byte(sendMail.Html)
	}
	//err := e.SendWithTLS(settings.Config.SMTP.Host+":"+strconv.Itoa(settings.Config.SMTP.Port), smtp.PlainAuth("", settings.Config.SMTP.User, settings.Config.SMTP.Pass, settings.Config.SMTP.Host), &tls.Config{InsecureSkipVerify: true})

	var err error

	if conf.User == "" {
		err = e.SendWithStartTLS(conf.Host+":"+strconv.Itoa(conf.Port), nil, &tls.Config{InsecureSkipVerify: true})
	} else {
		auth := smtp.PlainAuth("", conf.User, conf.Pass, conf.Host)
		// err := e.SendWithTLS(settings.Config.SMTP.Host+":"+strconv.Itoa(settings.Config.SMTP.Port), auth, &tls.Config{InsecureSkipVerify: true})
		err = e.Send(conf.Host+":"+strconv.Itoa(conf.Port), auth)
		//err := e.SendWithTLS(settings.Config.SMTP.Host+":"+strconv.Itoa(settings.Config.SMTP.Port), smtp.PlainAuth("", settings.Config.SMTP.User, settings.Config.SMTP.Pass, settings.Config.SMTP.Host), &tls.Config{InsecureSkipVerify: true})
	}

	return err

}
