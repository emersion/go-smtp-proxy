package proxy

import (
	"crypto/tls"
	"net"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

type Security int

const (
	SecurityNone Security = iota
	SecuritySTARTTLS
	SecurityTLS
)

type Backend struct {
	Addr string
	Security Security
	TLSConfig *tls.Config

	unexported struct{}
}

func New(addr string) *Backend {
	return &Backend{Addr: addr}
}

func NewTLS(addr string, tlsConfig *tls.Config) *Backend {
	return &Backend{
		Addr: addr,
		Security: SecurityTLS,
		TLSConfig: tlsConfig,
	}
}

func (be *Backend) login(username, password string) (*smtp.Client, error) {
	var conn net.Conn
	var err error
	if be.Security == SecurityTLS {
		conn, err = tls.Dial("tcp", be.Addr, be.TLSConfig)
	} else {
		conn, err = net.Dial("tcp", be.Addr)
	}
	if err != nil {
		return nil, err
	}

	host, _, _ := net.SplitHostPort(be.Addr)
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return nil, err
	}

	if be.Security == SecuritySTARTTLS {
		if err := c.StartTLS(be.TLSConfig); err != nil {
			return nil, err
		}
	}

	auth := sasl.NewPlainClient("", username, password)
	if err := c.Auth(auth); err != nil {
		return nil, err
	}

	return c, nil
}

func (be *Backend) Login(username, password string) (smtp.User, error) {
	c, err := be.login(username, password)
	if err != nil {
		return nil, err
	}

	u := &user{
		c: c,
		be: be,
		username: username,
	}
	return u, nil
}
