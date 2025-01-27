package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"rate_limiter/limiter"
	"testing"
	"time"
)

type rateLimiterTestCase struct {
	name          string
	config        Config
	totalRequests int
	expectCode    int
	apiKey        bool
}

func TestRateLimiter(t *testing.T) {
	var testCases = []rateLimiterTestCase{
		{
			name: "Allow requests within IP rate limit",
			config: Config{
				IPRateLimit:    5,
				TokenRateLimit: 4,
				IPBlockTime:    1,
				TokenBlockTime: 1,
			},
			totalRequests: 4,
			expectCode:    http.StatusOK,
			apiKey:        false,
		},
		{
			name: "Block requests exceeding IP rate limit",
			config: Config{
				IPRateLimit:    5,
				TokenRateLimit: 4,
				IPBlockTime:    1,
				TokenBlockTime: 1,
			},
			totalRequests: 10,
			expectCode:    http.StatusTooManyRequests,
			apiKey:        false,
		},
		{
			name: "Allow requests within token rate limit using API key",
			config: Config{
				IPRateLimit:    5,
				TokenRateLimit: 8,
				IPBlockTime:    1,
				TokenBlockTime: 1,
			},
			totalRequests: 7,
			expectCode:    http.StatusOK,
			apiKey:        true,
		},
		{
			name: "Block requests exceeding token rate limit using API key",
			config: Config{
				IPRateLimit:    5,
				TokenRateLimit: 8,
				IPBlockTime:    1,
				TokenBlockTime: 1,
			},
			totalRequests: 9,
			expectCode:    http.StatusTooManyRequests,
			apiKey:        true,
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			redisStore := limiter.NewRedisStore("localhost:6379", "")
			limiter.TestIP = fmt.Sprintf("192.168.1.%d", i+1)

			apiKey := ""
			if tc.apiKey {
				apiKey = fmt.Sprintf("test-api-key-%d", i+rand.Int())
			}

			rateLimiter := limiter.NewRateLimiter(
				redisStore,
				int64(tc.config.IPRateLimit),
				int64(tc.config.TokenRateLimit),
				time.Duration(tc.config.IPBlockTime)*time.Minute,
				time.Duration(tc.config.TokenBlockTime)*time.Minute,
			)

			server := httptest.NewServer(rateLimiter.MiddlewareHTTP(http.HandlerFunc(helloHandler)))
			defer server.Close()

			client := &http.Client{}

			var lastStatusCode int

			for i = 0; i < tc.totalRequests; i++ {
				req, err := http.NewRequest("GET", server.URL, nil)
				if tc.apiKey {
					req.Header.Set("API_KEY", apiKey)
				}
				resp, err := client.Do(req)
				if err != nil {
					t.Fatalf("Erro ao fazer requisição: %v", err)

				}
				lastStatusCode = resp.StatusCode
			}
			if lastStatusCode != tc.expectCode {
				t.Errorf("Esperando status %d, mas recebeu %d", tc.expectCode, lastStatusCode)
			}

		})
	}
}
