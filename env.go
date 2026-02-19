package main

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	// DBConnStr  string `env:"DB_CONN_STR,required"`
	VaultClientID string `env:"VAULT_CLIENT_ID"`
	VaultAddr     string `env:"VAULT_ADDR"`
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found or error loading it, fallback to lookup from os")
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Panic("failed to parse env:", err)
	}

	log.Printf("Config: %+v\n", cfg)

	return &cfg
}
