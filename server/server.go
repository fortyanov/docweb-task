package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"docweb-task/log"
)

var (
	config   *Config
	server   *http.Server
	ctx      context.Context
)

func Init(cfg *Config) (err error) {
	config = cfg

	router := httprouter.New()

	router.POST("/", Processing(uploadHandler, uploadPreCallback, uploadPostCallback))
	router.GET("/", Processing(downloadHandler, downloadPreCallback, downloadPostCallback))
	router.DELETE("/", Processing(deleteHandler, deletePreCallback, deletePostCallback))

	router.GlobalOPTIONS = http.HandlerFunc(optionsHandler)
	router.PanicHandler = panicHandler

	server = &http.Server{Addr: fmt.Sprintf("%s:%d", config.Host, config.Port), Handler: router}

	return
}

func Run(c context.Context, cancel context.CancelFunc) {
	ctx = c
	go func() {
		defer cancel()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server error: ", err)
		}
	}()
}

func Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Gracefully shutdown error: ", err)
	}
}
