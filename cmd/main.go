package main

import (
	"embed"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/maddsua/flippercardapp/justserve"
)

//go:embed web/*
var webfs embed.FS

func main() {

	mux := http.NewServeMux()

	mux.Handle("/api/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte("Not implemented. Yet."))
	}))

	mux.Handle("/", justserve.NewStaticAssetServer(webfs, "web/dist"))

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", envServePort(80)),
		Handler: mux,
	}

	exitCh := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		errChan <- srv.ListenAndServe()
	}()

	slog.Info("Starting http server",
		slog.String("addr", srv.Addr))

	select {
	case err := <-errChan:
		slog.Error("HTTP server: Listen",
			slog.String("err", err.Error()))
	case <-exitCh:
		if err := srv.Close(); err != nil {
			slog.Error("HTTP server: Close",
				slog.String("err", err.Error()))
		}
	}
}

func envServePort(defaultPort int) int {
	if port, _ := strconv.Atoi(os.Getenv("PORT")); port >= 80 && port < math.MaxUint16 {
		return port
	}
	return defaultPort
}
