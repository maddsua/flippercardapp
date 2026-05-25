package main

import (
	"embed"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	db_pkg "github.com/maddsua/flippercardapp/db"
	"github.com/maddsua/flippercardapp/media"
	"github.com/maddsua/flippercardapp/rest"
	"github.com/maddsua/flippercardapp/spa"
)

var Version = "development"
var BuildTS = "unknown"

//go:embed web/*
var WebFS embed.FS

func main() {

	godotenv.Load()

	if strings.EqualFold(os.Getenv("LOG_FMT"), "json") {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	}

	if strings.EqualFold(os.Getenv("DEBUG"), "true") {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("ENABLED")
	}

	slog.Info("FLIPPERCARD APP",
		slog.String("version", Version),
		slog.String("build_ts", BuildTS))

	dataDir := "data"
	if val := os.Getenv("DATA_DIR"); val != "" {
		dataDir = val
	}

	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		slog.Error("Prepare workdir",
			slog.String("err", err.Error()))
		os.Exit(1)
	}

	slog.Info("Workdir OK",
		slog.String("path", dataDir))

	dbconn, err := db_pkg.Open(path.Join(dataDir, "index.db3"))
	if err != nil {
		slog.Error("Open database",
			slog.String("err", err.Error()))
		os.Exit(1)
	}

	defer dbconn.Close()

	slog.Info("DB OK")

	if !strings.EqualFold(os.Getenv("DB_MIGRATE"), "false") {

		state, err := db_pkg.Migrate(dbconn)
		if err != nil {
			slog.Error("DB migration failed",
				slog.String("err", err.Error()))
			os.Exit(1)
		}

		if state.Updated {
			slog.Info("DB schema updated",
				slog.Int("version", int(state.Version)))
		} else {
			slog.Info("DB schema up to date")
		}
	}

	if err := db_pkg.InitDatabase(dbconn, db_pkg.StateInitParams{
		RootUserName:     os.Getenv("INIT_ROOT_USERNAME"),
		RootUserPassword: os.Getenv("INIT_ROOT_PASSWORD"),
		RestRootPassword: strings.EqualFold(os.Getenv("INIT_RESET_ROOT_PASSWORD"), "true"),
	}); err != nil {
		slog.Error("DB state init failed",
			slog.String("err", err.Error()))
		os.Exit(1)
	}

	mux := http.NewServeMux()

	mux.Handle("/media/", http.StripPrefix("/media", media.NewHandler(dbconn)))
	mux.Handle("/api/", http.StripPrefix("/api", rest.NewHandler(dbconn)))
	mux.Handle("/", spa.NewServerSPA(spa.NewEFSSnapshot(WebFS, BuildTS), "web/dist"))

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", envServePort(8280)),
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
