package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"distributed-job-scheduler/internal/config"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	conf, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	startServer(conf)

	<-ctx.Done()
	slog.Info("Shutting down application")
}

func startServer(conf config.Config) {
	slog.Info("Starting server on ", "addr", conf.ServerAddr, "port", conf.ServerPort)

	http.HandleFunc("/", indexHandler)
	go func() {
		addr := conf.ServerAddr + ":" + conf.ServerPort
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			slog.Error("Failed to start server", "err", err)
		}
	}()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Scheduler!")
}
