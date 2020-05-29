package storage

import (
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

var (
	config *Config
)

func Init(cfg *Config) (err error) {
	config = cfg
	return
}

func createFile(path string) (file *os.File, err error) {
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)

	file, err = os.OpenFile(filepath.Join(path), os.O_WRONLY|os.O_CREATE, 0666)
	return
}

func dirIsEmpty(path string) (isEmpty bool, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return
	}
	defer f.Close()

	if _, err = f.Readdirnames(1); err == io.EOF {
		isEmpty = true
		err = nil
	}

	return
}

func Upload(formFile multipart.File, formHeader *multipart.FileHeader) (fileHash string, err error) {
	var (
		file *os.File
		h    hash.Hash
	)

	h = sha1.New()
	if _, err = io.Copy(h, formFile); err != nil {
		return fileHash, fmt.Errorf("File '%s' creation hash error: %v", formHeader.Filename, err)
	}
	fileHash = fmt.Sprintf("%x", h.Sum(nil))

	path := filepath.Join(config.Directory, fileHash[:2], fileHash)
	if file, err = createFile(path); err != nil {
		return fileHash, fmt.Errorf("File '%s' creation error: %v", formHeader.Filename, err)
	}
	defer file.Close()

	formFile.Seek(0, 0)
	if _, err = io.Copy(file, formFile); err != nil {
		return fileHash, fmt.Errorf("File '%s' upload error: %v", formHeader.Filename, err)
	}

	return
}

func Delete(fileHash string) (err error) {
	path := filepath.Join(config.Directory, fileHash[:2], fileHash)
	if err = os.Remove(path); err != nil {
		return fmt.Errorf("File '%s' delete error: %v", fileHash, err)
	}

	dir := filepath.Dir(path)
	isEmpty, _ := dirIsEmpty(dir)
	if isEmpty == true {
		_ = os.Remove(dir)
	}

	return
}

func Download(fileHash string) (file *os.File, err error) {
	path := filepath.Join(config.Directory, fileHash[:2], fileHash)
	file, err = os.Open(path)
	return
}
