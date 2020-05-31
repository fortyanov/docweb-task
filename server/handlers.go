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

func processing(handle httprouter.Handle, preCallback, postCallback func()) httprouter.Handle {
	return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		preCallback()
		defer postCallback()
		handle(writer, request, params)
	}
}

func uploadHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var (
		err          error
		formHashType storage.HashType
		formFile     multipart.File
		formHeader   *multipart.FileHeader
		fileHash     string
	)

	if request.ContentLength < config.MinUploadSize || request.ContentLength > config.MaxUploadSize {
		JsonError(writer, errors.New(errRequestBodySize), http.StatusBadRequest)
		return
	}
	var ok bool
	if formHashType, ok = storage.MapHashType[request.FormValue("hash-type")]; !ok {
		JsonError(writer, errors.New(errHashType), http.StatusBadRequest)
		return
	}
	if formFile, formHeader, err = request.FormFile("file"); err != nil {
		log.Warning(err)
		JsonError(writer, errors.New(errFormFile), http.StatusBadRequest)
		return
	}
	defer formFile.Close()

	if fileHash, err = storage.Upload(formHashType, formFile, formHeader); err != nil {
		log.Warning(err)
		JsonError(writer, errors.New(errUpload), http.StatusBadRequest)
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
		JsonError(writer, errors.New(errHashFile), http.StatusBadRequest)
		return
	}

	if err = storage.Delete(fileHash); err != nil {
		log.Warning(err)
		JsonError(writer, errors.New(errDelete), http.StatusBadRequest)
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func downloadHandler(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	var (
		err          error
		fileHash     string
		fileHashType storage.HashType
		file         *os.File
	)

	urlQueryValues := request.URL.Query()
	fileHash = urlQueryValues.Get("hash")
	if fileHash == "" || !storage.IsValidHash(fileHash) {
		log.Warning(err)
		JsonError(writer, errors.New(errHashFile), http.StatusBadRequest)
		return
	}
	var ok bool
	if fileHashType, ok = storage.MapHashType[urlQueryValues.Get("hash-type")]; !ok {
		log.Warning(err)
		JsonError(writer, errors.New(errHashType), http.StatusBadRequest)
		return
	}

	if file, err = storage.Download(fileHash); err != nil {
		log.Warning(err)
		JsonError(writer, errors.New(errDownload), http.StatusBadRequest)
		return
	}
	defer file.Close()

	var calculatedFileHash string
	calculatedFileHash, err = storage.CalcFileHash(fileHashType, file)
	if fileHash != calculatedFileHash {
		log.Warning(err)
		JsonError(writer, errors.New(errCorruptedFile), http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileHash))
	if _, err = io.Copy(writer, file); err != nil {
		log.Warning(err)
		JsonError(writer, errors.New(errDownload), http.StatusBadRequest)
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
