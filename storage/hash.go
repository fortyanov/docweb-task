package storage

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"io"
)

type HashType int

const (
	Md5 HashType = iota
	Sha1
	Sha256
)

var MapHashType = map[string]HashType{
	"md5":    Md5,
	"sha1":   Sha1,
	"sha256": Sha256,
}

func (h HashType) String() string {
	return [...]string{"md5", "sha1", "sha256"}[h]
}

func CalcFileHash(hashType HashType, file io.Reader) (fileHash string, err error) {
	var h hash.Hash

	switch hashType {
	case Md5:
		h = md5.New()
	case Sha1:
		h = sha1.New()
	case Sha256:
		h = sha256.New()
	default:
		return fileHash, errors.New(errIncorrectHashType)
	}

	if _, err = io.Copy(h, file); err != nil {
		return fileHash, err
	}
	fileHash = fmt.Sprintf("%x", h.Sum(nil))

	return fileHash, err
}
