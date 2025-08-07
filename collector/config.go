package collector

import "errors"

type Config struct {
	Username string
	Password string
}

var (
	errNoUsername = errors.New("username must be specified")
	errNoAuth     = errors.New("password or private_key must be specified")
)

func (c Config) Validate() error {
	if c.Username == "" {
		return errNoUsername
	}

	if c.Password == "" {
		return errNoAuth
	}

	return nil
}
