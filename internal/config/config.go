package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

const (
	serverAddrDefault   = ":8080"
	baseShortURLDefault = "http://localhost:8080"
	fileStoragePath     = ""
)

type Config struct {
	ServerAddr      string
	BaseShortURL    string
	FileStoragePath string
}

func LoadConfig() (Config, error) {
	cfg := Config{
		ServerAddr:      serverAddrDefault,
		BaseShortURL:    baseShortURLDefault,
		FileStoragePath: fileStoragePath,
	}

	cfg.loadEnv()

	cfg.loadFlags()
	if err := cfg.validate(); err != nil {
		return cfg, fmt.Errorf("cant validate config: %w", err)
	}

	return cfg, nil
}

func (cfg *Config) loadEnv() {
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		cfg.ServerAddr = envRunAddr
	}

	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		cfg.BaseShortURL = envBaseAddr
	}

	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	}
}

func (cfg *Config) loadFlags() {
	flag.StringVar(&cfg.ServerAddr, "a", cfg.ServerAddr, "HTTP server address")
	flag.StringVar(&cfg.BaseShortURL, "b", cfg.BaseShortURL, "Base shortened URL")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "File storage path")

	flag.Parse()
}

func (cfg *Config) validate() error {
	_, err := url.Parse(cfg.BaseShortURL)
	if err != nil {
		return fmt.Errorf("cant parse base short ulr: %w", err)
	}

	return nil
}
