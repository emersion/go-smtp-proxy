package proxy

import (
	"io"
	"net/smtp"

	server "github.com/emersion/go-smtp-server"
)

type user struct {
	c *smtp.Client
	be *Backend
	username string
}

func (u *user) Send(msg *server.Message) error {
	if err := u.c.Mail(msg.From); err != nil {
		return err
	}
	for _, to := range msg.To {
		if err := u.c.Rcpt(to); err != nil {
			return err
		}
	}

	wc, err := u.c.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	_, err = io.Copy(wc, msg.Data)
	return err
}

func (u *user) Logout() error {
	return u.c.Close()
}
