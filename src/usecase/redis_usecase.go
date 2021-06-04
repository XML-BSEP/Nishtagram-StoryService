package usecase

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisUseCase interface {
	AddKeyValueSet(context context.Context, key string, value interface{},  expiration time.Duration) error
	GetValueByKey(context context.Context, key string) (string, error)
	DeleteValueByKey(context context.Context, key string) error
	ExistsByKey(context context.Context, key string) bool
	ScanKeyByPattern(context context.Context, pattern string) ([]string, error)
}

type redisUseCase struct {
	RedisClient *redis.Client
}

func (r *redisUseCase) ScanKeyByPattern(ctx context.Context, pattern string) ([]string, error) {
	var _, cursor uint64
	var n int
	var keysLength []string
	keysLength, _, _ = r.RedisClient.Scan(ctx, cursor, "*", 0).Result()

	numFor := len(keysLength)
	for {

		var keys []string
		var err error
		keys, cursor, err = r.RedisClient.Scan(context.Background(), cursor, pattern, int64(numFor)).Result()
		if err != nil {
			panic(err)
		}
		n += len(keys)
		if len(keys) > 0 {
			return keys, nil
		}

		if cursor == 0 {
			break
		}
	}

	return nil, fmt.Errorf("server error")
}

func NewRedisUsecase(r *redis.Client) RedisUseCase {
	return &redisUseCase{RedisClient: r}
}

func (r *redisUseCase) GetValueByKey(context context.Context, key string) (string, error) {
	return r.RedisClient.Get(context, key).Result()
}

func (r *redisUseCase) AddKeyValueSet(context context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.RedisClient.Set(context, key, value, expiration).Err()
}

func (r *redisUseCase) DeleteValueByKey(context context.Context, key string) error {
	return r.RedisClient.Del(context, key).Err()
}



func (r *redisUseCase) ExistsByKey(context context.Context, key string) bool {
	res := r.RedisClient.Exists(context, key).Val()
	if res == 0 {
		return false
	}
	return true
}
