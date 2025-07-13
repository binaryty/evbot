package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"

	"github.com/binaryty/evbot/internal/app"
	"github.com/binaryty/evbot/internal/config"
)

func main() {
	cfg := config.Load()
	app := app.NewApp(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := app.Init(ctx); err != nil {
		log.Fatalf("App init failed: %v", err)
	}

	if err := app.Run(ctx); err != nil {
		log.Fatalf("App runtime error: %v", err)
	}
}
