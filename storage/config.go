package storage

import (
	"errors"
)

type Config struct {
	Directory string `ini:"directory"`
	MaxSize   int64  `ini:"max_size"`
}

func (c *Config) Validate() error {
	if c.Directory == "" {
		return errors.New(errIncorrectStorageDir)
	}
	if c.MaxSize < 1000 {
		return errors.New(errIncorrectStorageMaxSize)
	}
	return nil
}
