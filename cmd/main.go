package main

import (
	"context"
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"os/signal"

	"github.com/binaryty/evbot/internal/config"
	"github.com/binaryty/evbot/internal/delivery/telegram"
	"github.com/binaryty/evbot/internal/repository/sqlite"
	"github.com/binaryty/evbot/internal/usecase"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := config.Load()

	// Инициализация БД
	db := initDB(cfg.DBPath)
	defer db.Close()

	bot := initBot(cfg.BotToken)

	eventRepo := sqlite.NewEventRepository(db)
	userRepo := sqlite.NewUserRepository(db)
	stateRepo := sqlite.NewStateRepository(db)
	registrationRepo := sqlite.NewRegistrationRepository(db)

	eventUC := usecase.NewEventUseCase(eventRepo)
	registrationUC := usecase.NewRegistrationUseCase(eventRepo, registrationRepo)
	userUC := usecase.NewUserUseCase(userRepo)

	handler := telegram.NewHandler(
		cfg,
		bot,
		eventUC,
		registrationUC,
		userUC,
		userRepo,
		stateRepo,
	)

	u := tgbotapi.NewUpdate(0)
	updates := bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			if err := handler.HandleUpdate(ctx, &update); err != nil {
				log.Printf("can't handle update: %v", err)
			}
		}
	}

}

func initDB(dsn string) *sql.DB {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}

	//if _, err := db.Exec(readMigrations()); err != nil {
	//	log.Fatal(err)
	//}

	return db
}

func initBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return bot
}
