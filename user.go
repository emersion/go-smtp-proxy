package proxy

import (
	"io"

	"github.com/emersion/go-smtp"
)

type user struct {
	c  *smtp.Client
	be *Backend
}

func (u *user) Send(from string, to []string, r io.Reader) error {
	if err := u.c.Mail(from); err != nil {
		return err
	}
	for _, rcpt := range to {
		if err := u.c.Rcpt(rcpt); err != nil {
			return err
		}
	}

	wc, err := u.c.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	_, err = io.Copy(wc, r)
	return err
}

func (u *user) Logout() error {
	return u.c.Close()
}
