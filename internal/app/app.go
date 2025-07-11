package app

import (
	"context"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"os"

	"github.com/binaryty/evbot/internal/config"
	"github.com/binaryty/evbot/internal/delivery/telegram"
	"github.com/binaryty/evbot/internal/repository/sqlite"
	"github.com/binaryty/evbot/internal/usecase"
)

type App struct {
	cfg    *config.Config
	logger *slog.Logger
}

func NewApp(
	cfg *config.Config,
	logger *slog.Logger,
) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

// Start ...
func (a *App) Start(ctx context.Context) {
	db := a.initDB()
	bot := a.initBot()

	logger := a.initLogger()

	eventRepo := sqlite.NewEventRepository(db)
	userRepo := sqlite.NewUserRepository(db)
	stateRepo := sqlite.NewStateRepository(db)
	registrationRepo := sqlite.NewRegistrationRepository(db)

	eventUC := usecase.NewEventUseCase(eventRepo)
	userUC := usecase.NewUserUseCase(userRepo)
	registrationUC := usecase.NewRegistrationUseCase(eventRepo, registrationRepo)

	handler := telegram.NewHandler(a.cfg, bot, logger, eventUC, registrationUC, userUC, stateRepo)

	u := tgbotapi.NewUpdate(0)
	updates := bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			if err := handler.HandleUpdate(ctx, &update); err != nil {
				a.logger.Error("can't handle update", slog.String("[error]", err.Error()))
			}
		}
	}
}

// initDB ...
func (a *App) initDB() *sql.DB {
	db, err := sql.Open("sqlite3", a.cfg.DBPath)
	if err != nil {
		panic("failed to init db " + err.Error())
	}

	return db
}

// initBot ...
func (a *App) initBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(a.cfg.BotToken)
	if err != nil {
		panic("failed to init bot " + err.Error())
	}
	a.logger.Info("bot successfully initialized",
		slog.String("UserName", bot.Self.UserName))

	return bot
}

// initLogger ...
func (a *App) initLogger() *slog.Logger {

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
