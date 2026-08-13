package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" default:":8080"`
	BaseURL       string `env:"BASE_URL" default:"http://localhost:8080"`
	LogLevel      string `env:"LOG_LEVEL" default:"info"`
}

func LoadFromCmd() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.ServerAddress, "a", ":8080", "ServerAddress to listen on")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "ServerAddress to listen on")
	flag.Parse()

	return cfg
}

func Load() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.ServerAddress, "a", ":8080", "ServerAddress to listen on")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "base URL for shortened links")
	flag.StringVar(&cfg.LogLevel, "l", "info", "log level (debug, info, warn, error)")
	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// if serverAddress, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
	// 	cfg.ServerAddress = serverAddress
	// }

	// if baseURL, ok := os.LookupEnv("BASE_URL"); ok {
	// 	cfg.BaseURL = baseURL
	// }

	return cfg
}
