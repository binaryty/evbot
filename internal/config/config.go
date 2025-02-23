package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	BotToken string
	DBPath   string
	AdminID  int64
}

func Load() *Config {
	admID, err := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	if err != nil {
		log.Fatalf("failed to parse admin ID")
	}

	botToken := os.Getenv("BOT_TOKEN")
	dbPath := os.Getenv("DB_PATH")

	if botToken == "" || dbPath == "" {
		log.Fatalf("environment variables must be specified")
	}

	return &Config{
		BotToken: botToken,
		DBPath:   dbPath,
		AdminID:  admID,
	}
}
