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
	"github.com/lindritprekaj/notification-service/internal/infrastructure/config"
	httpinfra "github.com/lindritprekaj/notification-service/internal/infrastructure/http"
	"github.com/lindritprekaj/notification-service/internal/infrastructure/notifier"
	"github.com/lindritprekaj/notification-service/internal/infrastructure/rabbitmq"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config error", "err", err)
		os.Exit(1)
	}

	mock := notifier.NewMock()
	dispatcher := application.NewDispatcher(mock.Send)

	consumer := rabbitmq.NewConsumer(cfg.RabbitMQURL, dispatcher.Dispatch)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Start RabbitMQ consumer in background (reconnects automatically)
	go consumer.Run(ctx)

	// HTTP health server
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
