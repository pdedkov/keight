package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"keight/internal/api"
	"keight/internal/app"
	"keight/internal/config"
	"keight/internal/db"
	"keight/internal/logging"
	"keight/internal/processor"
	"keight/internal/storage"

	"github.com/joho/godotenv"
)

const (
	readHeaderTimeout = 3 * time.Second
)

func main() {
	var err error
	// load env
	if os.Getenv("APP_MODE") == app.ModeDev {
		err = godotenv.Load(".env")
		if err != nil {
			slog.Error("cant load .env file", "error", err)
			os.Exit(1)
		}
	}

	// load app config
	var cfg *config.Config
	cfg, err = config.New()
	if err != nil {
		slog.Error("cant load config", "error", err)
		os.Exit(1)
	}
	// init logging
	log := logging.NewSLog(app.APIService, cfg.Log)

	redis, err := db.NewRedis(cfg.Redis, log)
	if err != nil {
		log.WithError(err).WithField(logging.Address, cfg.Redis.Address()).
			Fatal("cant connect to db redis")
	}
	router := http.NewServeMux()
	srv := &http.Server{
		Addr:              cfg.API.Addr(),
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	apiSrv := api.New(processor.New(storage.New(cfg.Storage.Count, cfg.Storage.ChunksCount, redis, log), log), log)
	apiSrv.ApplyHanders(router)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-ctx.Done()
		if err = srv.Shutdown(ctx); err != nil {
			log.WithError(err).Warn("cant graceful shutdown api service")
		}
		cancel()
	}()

	log.WithField(logging.Address, cfg.API.Addr()).Info("run api service")
	if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.WithError(err).WithField(logging.Address, cfg.API.Addr()).
			Error("cant run api service")
	}
	log.Info("api exit")
}
