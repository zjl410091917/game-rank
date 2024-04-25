package dao

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// GetCache 获取redis缓存
func GetCache(ctx context.Context, rd redis.UniversalClient, key string, data any) (ok bool, err error) {
	res, err := rd.Get(ctx, key).Result()
	if err == nil {
		err = json.Unmarshal([]byte(res), data)
	}
	ok = err == nil
	if !ok && errors.Is(err, redis.Nil) {
		err = nil
	}
	return
}

func SetCache(ctx context.Context, rd redis.UniversalClient, key string, value any) (err error) {
	cache, err := json.Marshal(value)
	if err != nil {
		return
	}
	err = rd.Set(ctx, key, cache, time.Minute*30).Err()
	return
}
