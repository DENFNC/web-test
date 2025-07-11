package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/DENFNC/web-test/config"
	"github.com/DENFNC/web-test/internal/app"
)

func main() {
	logger := initLogger()
	cfg := config.LoadConfig(logger, "./.env.example")

	application := app.NewApp(logger, cfg)
	go application.MustStart()

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh

	logger.Info(
		"Calling program termination",
		slog.String("signal", sig.String()),
	)
}

func initLogger() *slog.Logger {
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{},
		),
	)
	slog.SetDefault(logger)

	return logger
}
