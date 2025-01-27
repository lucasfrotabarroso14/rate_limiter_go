package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

type Config struct {
	IPRateLimit    int
	TokenRateLimit int
	IPBlockTime    int
	TokenBlockTime int
}

func loadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env")
	}
	getEnvInt := func(key string) int {
		value, err := strconv.Atoi(os.Getenv(key))
		if err != nil {
			log.Fatalf("Erro ao converter %s: %v", key, err)
		}
		return value
	}
	return Config{
		IPRateLimit:    getEnvInt("IP_RATE_LIMIT"),
		TokenRateLimit: getEnvInt("TOKEN_RATE_LIMIT"),
		IPBlockTime:    getEnvInt("IP_BLOCK_TIME"),
		TokenBlockTime: getEnvInt("TOKEN_BLOCK_TIME"),
	}
}
