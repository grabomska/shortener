package config

import "flag"

type Config struct {
	Address       string `default:":8080"`
	ResultAddress string `default:"http://localhost:8080"`
}

func LoadFromCmd() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.Address, "a", ":8080", "address to listen on")
	flag.StringVar(&cfg.ResultAddress, "b", "http://localhost:8080", "address to listen on")
	flag.Parse()

	return cfg
}
