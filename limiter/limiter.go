package limiter

import (
	"context"
	"net/http"
	"strings"
	"time"
)

var TestIP string

type Store interface {
	Increment(ctx context.Context, key string, expiration string) (int64, error)
	IsBlocked(ctx context.Context, key string) (bool, error)
	BlockKey(ctx context.Context, key string, blockTime time.Duration) error
}

type RateLimiter struct {
	store          *RedisStore
	ipRateLimit    int64
	tokenRateLimit int64
	ipBlockTime    time.Duration
	tokenBlockTime time.Duration
}

func NewRateLimiter(store *RedisStore, ipRateLimit, tokenRateLimit int64, ipBlockTime, tokenBlockTime time.Duration) *RateLimiter {
	return &RateLimiter{
		store:          store,
		ipRateLimit:    ipRateLimit,
		tokenRateLimit: tokenRateLimit,
		ipBlockTime:    ipBlockTime,
		tokenBlockTime: tokenBlockTime,
	}
}

func (r *RateLimiter) MiddlewareHTTP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		token := req.Header.Get("API_KEY")

		clientIP := getClientIP(req)

		var key string
		var limit int64
		var blockTime time.Duration

		if token != "" {
			key = "token" + token
			limit = r.tokenRateLimit
			blockTime = r.tokenBlockTime

		} else {
			key = "ip:" + clientIP
			limit = r.ipRateLimit
			blockTime = r.ipBlockTime
		}

		isBlocked, err := r.store.IsBlocked(ctx, key)
		if err != nil {
			http.Error(w, "Erro ao verificar bloqueio", http.StatusInternalServerError)
			return
		}
		if isBlocked {
			http.Error(w, "429 Too Many Requests - Você atingiu o limite de requisicao", http.StatusTooManyRequests)
			return
		}

		counter, err := r.store.Increment(ctx, key, 100*time.Second)

		if counter >= limit {
			err = r.store.BlockKey(ctx, key, blockTime)
			if err != nil {
				http.Error(w, "erro ao bloquear", http.StatusInternalServerError)
				return
			}
		}
		next.ServeHTTP(w, req)

	},
	)
}

func getClientIP(r *http.Request) string {
	if TestIP != "" {
		return TestIP
	}
	ipPort := r.RemoteAddr
	if strings.Contains(ipPort, ":") {
		return strings.Split(ipPort, ":")[0]
	}
	return ipPort
}
