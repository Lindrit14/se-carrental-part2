package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lindritprekaj/notification-service/internal/application"
	"github.com/lindritprekaj/notification-service/internal/application/handlers"
	"github.com/lindritprekaj/notification-service/internal/infrastructure/config"
	"github.com/lindritprekaj/notification-service/internal/infrastructure/email"
	httpinfra "github.com/lindritprekaj/notification-service/internal/infrastructure/http"
	"github.com/lindritprekaj/notification-service/internal/infrastructure/notifier"
	"github.com/lindritprekaj/notification-service/internal/infrastructure/rabbitmq"
	"github.com/lindritprekaj/notification-service/internal/infrastructure/repository/sqlite"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config error", "err", err)
		os.Exit(1)
	}

	userRepo, err := sqlite.Open(cfg.UserDBPath)
	if err != nil {
		slog.Error("sqlite open", "err", err, "path", cfg.UserDBPath)
		os.Exit(1)
	}
	defer userRepo.Close()
	slog.Info("user read model ready", "path", cfg.UserDBPath)

	renderer, err := email.NewRenderer()
	if err != nil {
		slog.Error("email renderer", "err", err)
		os.Exit(1)
	}

	notif, err := notifier.New(notifier.Config{
		Type:                      cfg.NotifierType,
		ACSEndpoint:               cfg.ACSEndpoint,
		ACSSenderAddress:          cfg.ACSSenderAddress,
		ACSAuthMode:               cfg.ACSAuthMode,
		ACSConnectionStringFile:   cfg.ACSConnectionStringFile,
		ACSConnectionStringInline: cfg.ACSConnectionStringInline,
		FromName:                  cfg.FromName,
	})
	if err != nil {
		slog.Error("notifier init", "err", err)
		os.Exit(1)
	}
	slog.Info("notifier ready", "type", cfg.NotifierType)

	deps := &handlers.Deps{
		Notifier:        notif,
		UserRepo:        userRepo,
		Renderer:        renderer,
		FromName:        cfg.FromName,
		FrontendBaseURL: cfg.FrontendBaseURL,
	}

	dispatcher := application.NewDispatcher(deps)
	consumer := rabbitmq.NewConsumer(cfg.RabbitMQURL, dispatcher.Dispatch)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go consumer.Run(ctx)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      httpinfra.NewRouter(consumer),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	go func() {
		slog.Info("http server listening", "port", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server error", "err", err)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down")

	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutCtx)
}
