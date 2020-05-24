package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"docweb-task/config"
	"docweb-task/log"
	"docweb-task/pidfile"
	"docweb-task/server"
	"docweb-task/storage"
)

const defaultConfig = "config.ini"

func main() {
	os.Exit(realMain())
}

func realMain() int {
	var (
		err     error
		cfgFile string
		cfg     *config.Config
	)

	if len(os.Args) < 2 || os.Args[1] == "" {
		cfgFile = defaultConfig
	} else {
		cfgFile = os.Args[1]
	}

	if cfg, err = config.Parse(cfgFile); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Load config error: %s", err)
		return 1
	}

	if err = log.Init(cfg.General.LogDest, cfg.General.LogTag, cfg.General.LogLevel); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Init logger error: %s", err)
		return 1
	}
	defer log.Close()

	if err = pidfile.Write(cfg.General.PidFile); err != nil {
		log.Error("Create PID file error: ", err)
		return 1
	}
	defer pidfile.Unlink(cfg.General.PidFile)

	if err = storage.Init(&cfg.Storage); err != nil {
		log.Error("Initialize storage error: ", err)
		return 1
	}

	if err = server.Init(&cfg.Server); err != nil {
		log.Error("Initialize http server error: ", err)
		return 1
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server.Run(ctx, cancel)
	defer server.Shutdown()

	signalEvents := make(chan os.Signal, 1)
	signal.Notify(signalEvents, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	select {
	case s := <-signalEvents:
		log.Info(fmt.Sprintf("Caught signal %v: terminating", s))
	case <-ctx.Done():
		return 1
	}

	return 0
}
