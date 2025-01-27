package main

import (
	"fmt"
	"log"
	"net/http"
	"rate_limiter/limiter"
	"time"
)

func main() {

	redisStore := limiter.NewRedisStore("localhost:6379", "")

	config := loadConfig()
	fmt.Println("Configurações carregadas:")
	fmt.Printf("IP_RATE_LIMIT: %d\n", config.IPRateLimit)
	fmt.Printf("TOKEN_RATE_LIMIT: %d\n", config.TokenRateLimit)
	fmt.Printf("IP_BLOCK_TIME: %d minutos\n", config.IPBlockTime)
	fmt.Printf("TOKEN_BLOCK_TIME: %d minutos\n", config.TokenBlockTime)

	raterLimite := limiter.NewRateLimiter(
		redisStore,
		int64(config.IPRateLimit),
		int64(config.TokenRateLimit),
		time.Duration(config.IPBlockTime)*time.Minute,
		time.Duration(config.TokenBlockTime)*time.Minute)

	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: raterLimite.MiddlewareHTTP(mux),
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("erro ao iniciar servidor")
	}

}

func helloHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("teste \n"))
}
