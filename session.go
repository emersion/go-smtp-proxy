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

func (s *session) Send(from string, to []string, r io.Reader) error {
	if err := s.Mail(from); err != nil {
		return err
	}

	for _, rcpt := range to {
		if err := s.Rcpt(rcpt); err != nil {
			return err
		}
	}

	return s.Data(r)
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
