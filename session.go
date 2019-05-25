package proxy

import (
	"io"

	"github.com/emersion/go-smtp"
)

type session struct {
	c  *smtp.Client
	be *Backend
}

func (s *session) Reset() {
	s.c.Reset()
}

func (s *session) Mail(from string) error {
	return s.c.Mail(from)
}

func (s *session) Rcpt(to string) error {
	return s.c.Rcpt(to)
}

func (s *session) Data(r io.Reader) error {
	wc, err := s.c.Data()
	if err != nil {
		return err
	}

	_, err = io.Copy(wc, r)
	if err != nil {
		wc.Close()
		return err
	}

	return wc.Close()
}

func (s *session) Logout() error {
	return s.c.Quit()
}
