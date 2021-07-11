package proxy

import (
	"crypto/tls"
	"errors"
	"net"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
)

type Security int

const (
	SecurityTLS Security = iota
	SecurityStartTLS
	SecurityNone
)

type Backend struct {
	Addr      string
	Security  Security
	TLSConfig *tls.Config
	LMTP      bool
	Host      string
	LocalName string

	unexported struct{}
}

func New(addr string) *Backend {
	return &Backend{Addr: addr, Security: SecurityStartTLS}
}

func NewTLS(addr string, tlsConfig *tls.Config) *Backend {
	return &Backend{
		Addr:      addr,
		Security:  SecurityTLS,
		TLSConfig: tlsConfig,
	}
}

func NewLMTP(addr string, host string) *Backend {
	return &Backend{
		Addr:     addr,
		Security: SecurityNone,
		LMTP:     true,
		Host:     host,
	}
}

func (be *Backend) newConn() (*smtp.Client, error) {
	var conn net.Conn
	var err error
	if be.LMTP {
		if be.Security != SecurityNone {
			return nil, errors.New("smtp-proxy: LMTP doesn't support TLS")
		}
		conn, err = net.Dial("unix", be.Addr)
	} else if be.Security == SecurityTLS {
		conn, err = tls.Dial("tcp", be.Addr, be.TLSConfig)
	} else {
		conn, err = net.Dial("tcp", be.Addr)
	}
	if err != nil {
		return nil, err
	}

	var c *smtp.Client
	if be.LMTP {
		c, err = smtp.NewClientLMTP(conn, be.Host)
	} else {
		host := be.Host
		if host == "" {
			host, _, _ = net.SplitHostPort(be.Addr)
		}
		c, err = smtp.NewClient(conn, host)
	}
	if err != nil {
		return nil, err
	}

	if be.LocalName != "" {
		err = c.Hello(be.LocalName)
		if err != nil {
			return nil, err
		}
	}

	if be.Security == SecurityStartTLS {
		if err := c.StartTLS(be.TLSConfig); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (be *Backend) login(username, password string) (*smtp.Client, error) {
	c, err := be.newConn()
	if err != nil {
		return nil, err
	}

	auth := sasl.NewPlainClient("", username, password)
	if err := c.Auth(auth); err != nil {
		return nil, err
	}

	return c, nil
}

func (be *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	c, err := be.login(username, password)
	if err != nil {
		return nil, err
	}

	s := &session{
		c:  c,
		be: be,
	}
	return s, nil
}

func (be *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	c, err := be.newConn()
	if err != nil {
		return nil, err
	}

	s := &session{
		c:  c,
		be: be,
	}
	return s, nil
}
