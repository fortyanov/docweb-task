package server

import "net/http"

func corsHeaders(writer http.ResponseWriter, request *http.Request) {
	header := writer.Header()
	origin := request.Header.Get("Origin")
	if origin != "" {
		header.Set("Access-Control-Allow-Origin", origin)
		header.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		header.Set("Access-Control-Allow-Credentials", "true")
	}
}

func corsOptionHeaders(writer http.ResponseWriter, request *http.Request) {
	header := writer.Header()
	origin := request.Header.Get("Origin")
	if origin == "" {
		origin = "*"
	}
	header.Set("Access-Control-Allow-Origin", origin)
	header.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	header.Set("Access-Control-Allow-Headers", request.Header.Get("Access-Control-Request-Headers"))
	header.Set("Access-Control-Allow-Credentials", "true")
	header.Set("Access-Control-Max-Age", "1728000")
}
