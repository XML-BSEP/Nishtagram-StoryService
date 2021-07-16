package usecase

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	logger "github.com/jelena-vlajkov/logger/logger"
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
	logger *logger.Logger
}

func (r *redisUseCase) ScanKeyByPattern(ctx context.Context, pattern string) ([]string, error) {
	r.logger.Logger.Infof("scanning redis for stories to show by pattern %v\n", pattern)
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
	r.logger.Logger.Warnf("no stories to show by pattern %v\n", pattern)
	return nil, fmt.Errorf("server error")
}

func NewRedisUsecase(r *redis.Client, logger *logger.Logger) RedisUseCase {
	return &redisUseCase{RedisClient: r, logger: logger}
}

func (r *redisUseCase) GetValueByKey(context context.Context, key string) (string, error) {
	return r.RedisClient.Get(context, key).Result()
}

func (r *redisUseCase) AddKeyValueSet(context context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.RedisClient.Set(context, key, value, expiration).Err()
	if err != nil {
		r.logger.Logger.Errorf("error while adding key %v and value %v in redis, error: %v\n", key, value, err)
	}
	return err
}

func (r *redisUseCase) DeleteValueByKey(context context.Context, key string) error {
	err := r.RedisClient.Del(context, key).Err()
	if err != nil {
		r.logger.Logger.Errorf("error while deleting key %v in redis, error: %v\n", key, err)
	}
	return err
}



func (r *redisUseCase) ExistsByKey(context context.Context, key string) bool {
	res := r.RedisClient.Exists(context, key).Val()
	if res == 0 {
		return false
	}
	return true
}
