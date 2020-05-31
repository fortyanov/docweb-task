package storage

import "fmt"

const (
	errIncorrectHashType       = "Unsupported hash type."
	errIncorrectStorageDir     = "Incorrect storage directory."
	errIncorrectStorageMaxSize = "Incorrect storage maximum available size."
)

type UploadError struct {
	err      error
	fileHash string
	filename string
}

func (e *UploadError) Error() string {
	return fmt.Sprintf("File '%s' '%s' upload error: %v", e.fileHash, e.filename, e.err)
}

type DeleteError struct {
	err      error
	fileHash string
}

func (e *DeleteError) Error() string {
	return fmt.Sprintf("File '%s' delete error: %v", e.fileHash, e.err)
}

type DownloadError struct {
	err      error
	fileHash string
}

func (e *DownloadError) Error() string {
	return fmt.Sprintf("File '%s' download error: %v", e.fileHash, e.err)
}
