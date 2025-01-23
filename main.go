package main

import (
	"log"
	"net/http"
	"rate_limiter/limiter"
	"time"
)

//o servidor recebe uma requisicao
//middleware intercepta a requisicao
// o midleware -> verifica o ip ou token vendo se esta bloqueado no redis
// se sim -> retorna 429 too many requests
// se nao incrementa o contador no redis
// se o numero de req for maior que o limite grava um bloquio temporario no redis
// se tiver dentro do limite a requisicao pode passar para o handler final

func main() {

	redisStore := limiter.NewRedisStore("localhost:6379", "")

	raterLimite := limiter.NewRateLimiter(redisStore, 10, 2, 5*time.Minute, 3*time.Minute)

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
