package storage

import (
	"errors"
)

type Config struct {
	Directory string `ini:"directory"`
}

var (
	errDirectory = errors.New("incorrect storage directory")
)

func (c *Config) Validate() error {
	if c.Directory == "" {
		return errDirectory
	}
	return nil
}
