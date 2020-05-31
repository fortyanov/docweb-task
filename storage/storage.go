package storage

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"docweb-task/log"
)

var (
	config *Config
)

func Init(ctx context.Context, cfg *Config) (err error) {
	config = cfg
	go runCleaner(ctx)
	return err
}

func runCleaner(ctx context.Context) {
	var (
		size int64
		err  error
	)

	ticker := time.NewTicker(time.Duration(30) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			if size, err = dirSize(config.Directory); err != nil {
				log.Warning("Storage directory unavailable.")
			}
			if size > config.MaxSize {
				log.Info("Cleaning")
				// здесь должна быть реализация удаления редко используемых файлов,
				//но я не успеваю сделать часть с редисом поэтому пусть будет так
			}
		}
	}
}

func dirSize(path string) (size int64, err error) {
	err = filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func dirIsEmpty(path string) (isEmpty bool, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return isEmpty, err
	}
	defer f.Close()

	if _, err = f.Readdirnames(1); err == io.EOF {
		isEmpty = true
		err = nil
	}

	return isEmpty, err
}

func Upload(formHashType HashType, formFile multipart.File, formHeader *multipart.FileHeader) (fileHash string, err error) {
	var (
		file       *os.File
		h          hash.Hash
		fileReader io.Reader
	)

	switch formHashType {
	case Md5:
		h = md5.New()
	case Sha1:
		h = sha1.New()
	case Sha256:
		h = sha256.New()
	default:
		return fileHash, errors.New(errIncorrectHashType)
	}

	fileReader = io.TeeReader(formFile, h)
	if file, err = ioutil.TempFile("/tmp", "storage-"); err != nil {
		return fileHash, err
	}
	defer os.Remove(file.Name())
	defer file.Close()

	if _, err = io.Copy(file, fileReader); err != nil {
		return fileHash, &UploadError{fileHash: fileHash, filename: formHeader.Filename, err: err}
	}
	fileHash = fmt.Sprintf("%x", h.Sum(nil))
	path := filepath.Join(config.Directory, fileHash[:2], fileHash)
	dir := filepath.Dir(path)
	os.MkdirAll(dir, 0755)
	if err = os.Rename(file.Name(), path); err != nil {
		return fileHash, &UploadError{fileHash: fileHash, filename: formHeader.Filename, err: err}
	}

	return fileHash, err
}

func Delete(fileHash string) (err error) {
	path := filepath.Join(config.Directory, fileHash[:2], fileHash)
	if err = os.Remove(path); err != nil {
		return &DeleteError{fileHash: fileHash, err: err}
	}

	dir := filepath.Dir(path)
	isEmpty, _ := dirIsEmpty(dir)
	if isEmpty == true {
		_ = os.Remove(dir)
	}

	return err
}

func Download(fileHash string) (file *os.File, err error) {
	path := filepath.Join(config.Directory, fileHash[:2], fileHash)
	if file, err = os.Open(path); err != nil {
		return file, &DownloadError{fileHash: fileHash, err: err}
	}

	return file, err
}
