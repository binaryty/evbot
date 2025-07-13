package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/binaryty/evbot/internal/config"
	"github.com/binaryty/evbot/internal/delivery/telegram"
	"github.com/binaryty/evbot/internal/repository/sqlite"
	"github.com/binaryty/evbot/internal/usecase"
)

type App struct {
	cfg    *config.Config
	logger *slog.Logger
	db     *sql.DB
	bot    *tgbotapi.BotAPI
}

// NewApp ...
func NewApp(cfg *config.Config) *App {
	app := &App{
		cfg: cfg,
	}
	app.logger = app.initLogger()

	return app
}

// Init ...
func (a *App) Init(ctx context.Context) error {
	db, err := a.initDB()
	if err != nil {
		return fmt.Errorf("DB init failed", slog.String("error", err.Error()))
	}
	a.db = db

	bot, err := a.initBot()
	if err != nil {
		return fmt.Errorf("failed to init bot", slog.String("error", err.Error()))
	}
	a.bot = bot

	return nil
}

// Run ...
func (a *App) Run(ctx context.Context) error {
	defer func() {
		if a.db != nil {
			_ = a.db.Close()
		}
	}()

	eventRepo := sqlite.NewEventRepository(a.db)
	userRepo := sqlite.NewUserRepository(a.db)
	stateRepo := sqlite.NewStateRepository(a.db)
	registrationRepo := sqlite.NewRegistrationRepository(a.db)

	eventUC := usecase.NewEventUseCase(eventRepo, &config.Config{
		AdminIDs: a.cfg.AdminIDs,
	})
	userUC := usecase.NewUserUseCase(userRepo)
	registrationUC := usecase.NewRegistrationUseCase(eventRepo, registrationRepo)

	handler := telegram.NewHandler(a.cfg, a.bot, a.logger, eventUC, registrationUC, userUC, stateRepo)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30
	updates := a.bot.GetUpdatesChan(u)

	a.logger.Info("Bot is running...")

	for {
		select {
		case update := <-updates:
			if err := handler.HandleUpdate(ctx, &update); err != nil {
				a.logger.Error("can't handle update", slog.String("[error]", err.Error()))
			}
		case <-ctx.Done():
			a.logger.Info("shutting down bot...")
			return nil
		}
	}
}

// initDB ...
func (a *App) initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", a.cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init db: %w", err)
	}

	return db, nil
}

// initBot ...
func (a *App) initBot() (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(a.cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("failed to init bot: %w", err)
	}
	a.logger.Info("bot successfully initialized",
		slog.String("UserName", bot.Self.UserName))

	return bot, nil
}

// initLogger ...
func (a *App) initLogger() *slog.Logger {

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
