package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
)

func TestInitRedis(t *testing.T) {
	ctx := context.Background()
	InitRedis(&RedisConfig{
		Host:       "54.151.143.188",
		Port:       "6379",
		Username:   "pay",
		Auth:       "ZsB3DgltDHF8lz0O",
		DB:         0,
		PoolSize:   10,
		Encryption: 0,
		Framework:  "cluster",
		Prefix:     "pay:",
	})
	key := "pay:chain:pool:ETH:"
	members := []*redis.Z{
		&redis.Z{Member: "0xa9E1409ade429aEbcE0E58B321f8592c545B5066", Score: 1.211111},
		&redis.Z{Member: "0xa9E1409ade429aEbcE0E58B321f8592c545B5067", Score: 1.211222},
		&redis.Z{Member: "0xa9E1409ade429aEbcE0E58B321f8592c545B5065", Score: 1.211233},
	}
	for {
		fmt.Println(Redis.ZCard(ctx, key).Result())
		mem := Redis.ZPopMax(ctx, key, 1).Val()
		fmt.Println(mem)
		if len(mem) != 1 {
			fmt.Println(Redis.ZAdd(ctx, key, members...).Result())
			continue
		}
	}
}
