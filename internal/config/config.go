package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	BotToken string  `mapstructure:"bot_token"`
	DBPath   string  `mapstructure:"db_path"`
	AdminIDs []int64 `mapstructure:"admin_ids"`
}

func Load() *Config {
	admIDs := make([]int64, 0)
	id, err := strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	if err != nil {
		log.Fatalf("failed to parse admin ID")
	}
	admIDs = append(admIDs, id)

	botToken := os.Getenv("BOT_TOKEN")
	dbPath := os.Getenv("DB_PATH")

	if botToken == "" || dbPath == "" {
		log.Fatalf("environment variables must be specified")
	}

	return &Config{
		BotToken: botToken,
		DBPath:   dbPath,
		AdminIDs: admIDs,
	}
}
