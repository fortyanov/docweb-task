package server

import (
	"errors"
)

type Config struct {
	Host            string `ini:"host"`
	Port            int    `ini:"port"`
	MinUploadSize   int64  `ini:"min_upload_size"`
	MaxUploadSize   int64  `ini:"max_upload_size"`
	JwtPublicKey    string `ini:"jwt_public_key"`
	EncodingTimeout int    `ini:"encoding_timeout"`
}

var (
	errHost            = errors.New("server 'host' not defined")
	errPort            = errors.New("server 'port' invalid")
	errMinUploadSize   = errors.New("server 'min_upload_size' invalid value")
	errMaxUploadSize   = errors.New("server 'max_upload_size' invalid value")
)

func (c *Config) Validate() error {
	const (
		MinFileSize int64 = 1
		MaxFileSize int64 = 20 * 1024 * 1024
	)
	if c.Host == "" {
		return errHost
	}
	if c.Port < 0 || c.Port > 65536 {
		return errPort
	}
	if c.MinUploadSize > MaxFileSize || c.MinUploadSize < MinFileSize {
		return errMinUploadSize
	}
	if c.MaxUploadSize > MaxFileSize || c.MaxUploadSize < MinFileSize {
		return errMaxUploadSize
	}
	return nil
}
