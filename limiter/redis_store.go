package limiter

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(addr, password string) *RedisStore {

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return &RedisStore{
		client: rdb,
	}
}

func (r *RedisStore) Increment(ctx context.Context, key string, expiration time.Duration) (int64, error) {
	pipe := r.client.TxPipeline()

	counter := pipe.Incr(ctx, key)

	pipe.Expire(ctx, key, expiration)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return counter.Val(), nil
}

func (r *RedisStore) IsBlocked(ctx context.Context, key string) (bool, error) {
	blockKey := "block:" + key
	val, err := r.client.Get(ctx, blockKey).Result()

	// quando o usuário nao existe ou seja, nao esta bloqueado
	if err == redis.Nil {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	// se eu o valor da key que estou buscando é 1 , eu tenho que bloquear
	return val == "1", nil

}

func (r *RedisStore) BlockKey(ctx context.Context, key string, blockTime time.Duration) error {
	blockKey := "block:" + key
	err := r.client.Set(ctx, blockKey, "1", blockTime).Err()
	return err
}
