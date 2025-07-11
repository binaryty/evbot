package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	BotToken string  `yaml:"bot_token" env-required:"true"`
	DBPath   string  `yaml:"db_path" env-required:"true"`
	AdminIDs []int64 `yaml:"admin_ids"`
}

// Load ...
func Load() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("path to config file is not set")
	}

	return loadFromPath(path)
}

// loadFromPath ...
func loadFromPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath ...
func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config  file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}
	return path
}
