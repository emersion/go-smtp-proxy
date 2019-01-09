package proxy

import (
	"github.com/emersion/go-smtp"
)

var _ smtp.Backend = &Backend{}
