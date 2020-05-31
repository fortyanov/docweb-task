package server

const (
	errServerHost          = "Server 'host' not defined."
	errServerPort          = "Server 'port' invalid."
	errServerMinUploadSize = "Server 'min_upload_size' invalid value."
	errServerMaxUploadSize = "Server 'max_upload_size' invalid value."
	errHashType            = "Incorrect request hash-type."
	errFormFile            = "Incorrect request file."
	errRequestBodySize     = "Too large body."
	errUpload              = "Uploading error."
	errHashFile            = "Incorrect get parameter hash."
	errDelete              = "Deleting error."
	errDownload            = "Downloading error."
	errCorruptedFile       = "Corrupted file."
)
