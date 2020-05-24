package server

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"

	"docweb-task/log"
	"docweb-task/storage"
)

var (
	errFormFile        = errors.New("incorrect request file")
	errRequestBodySize = errors.New("too large body")
	errUpload          = errors.New("uploading error")
	errHashFile        = errors.New("incorrect get parameter hash")
	errDelete          = errors.New("deleting error")
	errDownload        = errors.New("downloading error")
)

func Processing(handle httprouter.Handle, preCallback, postCallback func()) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		preCallback()
		defer postCallback()
		handle(writer, request, params)
	}
}

func uploadHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var (
		err        error
		formFile   multipart.File
		formHeader *multipart.FileHeader
		fileHash   string
	)

	if request.ContentLength < config.MinUploadSize || request.ContentLength > config.MaxUploadSize {
		JsonError(writer, errRequestBodySize, http.StatusBadRequest)
		return
	}
	if formFile, formHeader, err = request.FormFile("file"); err != nil {
		log.Warning(err)
		JsonError(writer, errFormFile, http.StatusBadRequest)
		return
	}
	defer formFile.Close()

	if fileHash, err = storage.Upload(formFile, formHeader); err != nil {
		log.Warning(err)
		JsonError(writer, errUpload, http.StatusBadRequest)
		return
	}

	JsonFileHash(writer, fileHash, http.StatusCreated)
}

func deleteHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var (
		err      error
		fileHash string
	)

	urlQueryValues := request.URL.Query()
	fileHash = urlQueryValues.Get("hash")
	if fileHash == "" {
		log.Warning(err)
		JsonError(writer, errHashFile, http.StatusBadRequest)
		return
	}

	if err = storage.Delete(fileHash); err != nil {
		log.Warning(err)
		JsonError(writer, errDelete, http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func downloadHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var (
		err      error
		fileHash string
		file     *os.File
	)

	urlQueryValues := request.URL.Query()
	fileHash = urlQueryValues.Get("hash")
	if fileHash == "" {
		log.Warning(err)
		JsonError(writer, errHashFile, http.StatusBadRequest)
		return
	}

	if file, err = storage.Download(fileHash); err != nil {
		log.Warning(err)
		JsonError(writer, errDownload, http.StatusBadRequest)
		return
	}
	defer file.Close()

	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileHash))
	if _, err = io.Copy(writer, file); err != nil {
		log.Warning(err)
		JsonError(writer, errDownload, http.StatusBadRequest)
	}

	return
}

func optionsHandler(writer http.ResponseWriter, request *http.Request) {
	corsOptionHeaders(writer, request)
	writer.WriteHeader(http.StatusNoContent)
}

func panicHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {
	corsHeaders(writer, request)
	log.Error(err)
	JsonError(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
