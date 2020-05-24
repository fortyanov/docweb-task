package server

import (
	"fmt"
	"net/http"
)

func JsonError(writer http.ResponseWriter, err interface{}, code int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	fmt.Fprintf(writer, `{"error":"%v"}`, err)
}

func JsonFileHash(writer http.ResponseWriter, hash string, code int) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	fmt.Fprintf(writer, `{"hash":"%s"}`, hash)
}