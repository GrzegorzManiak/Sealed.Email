package server

import "github.com/emersion/go-smtp"

// The Backend implements SMTP server methods.
type Backend struct {
	Mode Mode
}

// NewSession is called after client greeting (EHLO, HELO).
func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}
