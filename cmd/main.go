package main

import (
	"context"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
	"os"
	"os/signal"

	"github.com/binaryty/evbot/internal/app"
	"github.com/binaryty/evbot/internal/config"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	App := app.NewApp(cfg, logger)

	App.Start(ctx)
}
