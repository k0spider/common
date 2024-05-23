package redis

import (
	"context"
	"crypto/tls"
	"github.com/go-redis/redis/v8"
	"github.com/k0spider/common/utils"
	"strings"
	"time"
)

type RedisConfig struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Username   string `yaml:"userName"`
	Auth       string `yaml:"auth"`
	DB         int    `yaml:"db"`
	PoolSize   int    `yaml:"poolSize"`
	Encryption uint8  `yaml:"encryption"`
	Framework  string `yaml:"framework"`
	Prefix     string `yaml:"prefix"`
	Tls        bool   `yaml:"tls"`
	MinVersion uint16 `yaml:"minVersion"`
}

var Redis *RedisClient

func InitRedis(config *RedisConfig) {
	auth := config.Auth
	if config.Encryption == 1 {
		auth = utils.MD5(auth)
	}
	addr := config.Host + config.Port
	if !strings.Contains(addr, ":") {
		addr = config.Host + ":" + config.Port
	}
	Redis = new(RedisClient)
	Redis.config = config
	if Redis.config.Framework == "cluster" {
		clusterOptions := &redis.ClusterOptions{Addrs: []string{addr}, PoolSize: config.PoolSize}
		if auth != "" {
			clusterOptions.DialTimeout = time.Second * 10
			clusterOptions.Username = config.Username
			clusterOptions.Password = auth
		}
		if config.Tls {
			if config.MinVersion == 0 {
				config.MinVersion = tls.VersionTLS12
			}
			clusterOptions.TLSConfig = &tls.Config{
				MinVersion: config.MinVersion,
			}
		}
		Redis.ClusterClient = redis.NewClusterClient(clusterOptions)
		if _, err := Redis.ClusterClient.Ping(context.Background()).Result(); err != nil {
			panic(err)
		}
	} else {
		options := &redis.Options{Addr: addr, PoolSize: config.PoolSize, DB: config.DB}
		if auth != "" {
			options.Username = config.Username
			options.Password = auth
		}
		if config.Tls {
			if config.MinVersion == 0 {
				config.MinVersion = tls.VersionTLS12
			}
			options.TLSConfig = &tls.Config{
				MinVersion: config.MinVersion,
			}
		}
		Redis.Client = redis.NewClient(options)
		if _, err := Redis.Client.Ping(context.TODO()).Result(); err != nil {
			panic(err)
		}
	}
}

type RedisClient struct {
	ClusterClient *redis.ClusterClient
	Client        *redis.Client
	config        *RedisConfig
}

func (r *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.Get(ctx, key)
	}
	return r.Client.Get(ctx, key)
}
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.Set(ctx, key, value, expiration)
	}
	return r.Client.Set(ctx, key, value, expiration)
}
func (r *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.SetNX(ctx, key, value, expiration)
	}
	return r.Client.SetNX(ctx, key, value, expiration)
}

func (r *RedisClient) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.SetEX(ctx, key, value, expiration)
	}
	return r.Client.SetEX(ctx, key, value, expiration)
}

func (r *RedisClient) HMSet(ctx context.Context, key string, values ...interface{}) *redis.BoolCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.HMSet(ctx, key, values...)
	}
	return r.Client.HMSet(ctx, key, values...)
}

func (r *RedisClient) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.HGetAll(ctx, key)
	}
	return r.Client.HGetAll(ctx, key)
}

func (r *RedisClient) HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.HSet(ctx, key, values...)
	}
	return r.Client.HSet(ctx, key, values...)
}

func (r *RedisClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.HGet(ctx, key, field)
	}
	return r.Client.HGet(ctx, key, field)
}

func (r *RedisClient) HIncrBy(ctx context.Context, key, field string, incr int64) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.HIncrBy(ctx, key, field, incr)
	}
	return r.Client.HIncrBy(ctx, key, field, incr)
}

func (r *RedisClient) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.HDel(ctx, key, fields...)
	}
	return r.Client.HDel(ctx, key, fields...)
}

func (r *RedisClient) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.HMGet(ctx, key, fields...)
	}
	return r.Client.HMGet(ctx, key, fields...)
}

func (r *RedisClient) HExists(ctx context.Context, key, field string) *redis.BoolCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.HExists(ctx, key, field)
	}
	return r.Client.HExists(ctx, key, field)
}

func (r *RedisClient) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.HKeys(ctx, key)
	}
	return r.Client.HKeys(ctx, key)
}

func (r *RedisClient) HLen(ctx context.Context, key string) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.HLen(ctx, key)
	}
	return r.Client.HLen(ctx, key)
}

func (r *RedisClient) HSetNX(ctx context.Context, key, field string, value interface{}) *redis.BoolCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.HSetNX(ctx, key, field, value)
	}
	return r.Client.HSetNX(ctx, key, field, value)
}

func (r *RedisClient) HRandField(ctx context.Context, key string, count int, withValues bool) *redis.StringSliceCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.HRandField(ctx, key, count, withValues)
	}
	return r.Client.HRandField(ctx, key, count, withValues)
}

func (r *RedisClient) HIncrByFloat(ctx context.Context, key, field string, incr float64) *redis.FloatCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.HIncrByFloat(ctx, key, field, incr)
	}
	return r.Client.HIncrByFloat(ctx, key, field, incr)
}

func (r *RedisClient) HVals(ctx context.Context, key string) *redis.StringSliceCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.HVals(ctx, key)
	}
	return r.Client.HVals(ctx, key)
}

func (r *RedisClient) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.HScan(ctx, key, cursor, match, count)
	}
	return r.Client.HScan(ctx, key, cursor, match, count)
}

func (r *RedisClient) ExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.ExpireAt(ctx, key, tm)
	}
	return r.Client.ExpireAt(ctx, key, tm)
}

func (r *RedisClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.Incr(ctx, key)
	}
	return r.Client.Incr(ctx, key)
}

func (r *RedisClient) Decr(ctx context.Context, key string) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.Decr(ctx, key)
	}
	return r.Client.Decr(ctx, key)
}

func (r *RedisClient) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.IncrBy(ctx, key, value)
	}
	return r.Client.IncrBy(ctx, key, value)
}

func (r *RedisClient) DecrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.DecrBy(ctx, key, value)
	}
	return r.Client.DecrBy(ctx, key, value)
}

func (r *RedisClient) IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.IncrByFloat(ctx, key, value)
	}
	return r.Client.IncrByFloat(ctx, key, value)
}

func (r *RedisClient) Del(ctx context.Context, key ...string) *redis.IntCmd {
	for k, v := range key {
		key[k] = r.config.Prefix + v
	}
	if r.config.Framework == "cluster" {
		return r.ClusterClient.Del(ctx, key...)
	}
	return r.Client.Del(ctx, key...)
}
func (r *RedisClient) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	pattern = r.config.Prefix + pattern
	if r.config.Framework == "cluster" {
		return r.ClusterClient.Keys(ctx, pattern)
	}
	return r.Client.Keys(ctx, pattern)
}

func (r *RedisClient) ZAddArgs(ctx context.Context, key string, args redis.ZAddArgs) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.ZAddArgs(ctx, key, args)
	}
	return r.Client.ZAddArgs(ctx, key, args)
}

func (r *RedisClient) ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.ZAdd(ctx, key, members...)
	}
	return r.Client.ZAdd(ctx, key, members...)
}
func (r *RedisClient) ZCard(ctx context.Context, key string) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.ZCard(ctx, key)
	}
	return r.Client.ZCard(ctx, key)
}

func (r *RedisClient) ZPopMax(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.ZPopMax(ctx, key, count...)
	}
	return r.Client.ZPopMax(ctx, key, count...)
}

func (r *RedisClient) ZRank(ctx context.Context, key string, member string) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.ZRank(ctx, key, member)
	}
	return r.Client.ZRank(ctx, key, member)
}

func (r *RedisClient) ZRevRank(ctx context.Context, key string, member string) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.ZRevRank(ctx, key, member)
	}
	return r.Client.ZRevRank(ctx, key, member)
}

func (r *RedisClient) ZScore(ctx context.Context, key string, member string) *redis.FloatCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.ZScore(ctx, key, member)
	}
	return r.Client.ZScore(ctx, key, member)
}
func (r *RedisClient) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.ZRangeWithScores(ctx, key, start, stop)
	}
	return r.Client.ZRangeWithScores(ctx, key, start, stop)
}
func (r *RedisClient) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.ZRevRangeWithScores(ctx, key, start, stop)
	}
	return r.Client.ZRevRangeWithScores(ctx, key, start, stop)
}

func (r *RedisClient) ZUnionStore(ctx context.Context, dest string, store *redis.ZStore) *redis.IntCmd {
	if r.config.Framework == "cluster" {
		return r.ClusterClient.ZUnionStore(ctx, dest, store)
	}
	return r.Client.ZUnionStore(ctx, dest, store)
}

func (r *RedisClient) ZUnionWithScores(ctx context.Context, store redis.ZStore) *redis.ZSliceCmd {
	if r.config.Framework == "cluster" {
		return r.ClusterClient.ZUnionWithScores(ctx, store)
	}
	return r.Client.ZUnionWithScores(ctx, store)
}

func (r *RedisClient) ZPopMin(ctx context.Context, key string, count ...int64) *redis.ZSliceCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.ZPopMin(ctx, key, count...)
	}
	return r.Client.ZPopMin(ctx, key, count...)
}

func (r *RedisClient) TTL(ctx context.Context, key string, count ...int64) *redis.DurationCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.TTL(ctx, key)
	}
	return r.Client.TTL(ctx, key)
}

func (r *RedisClient) Exists(ctx context.Context, key string) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.Exists(ctx, key)
	}
	return r.Client.Exists(ctx, key)
}

func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.Expire(ctx, key, expiration)
	}
	return r.Client.Expire(ctx, key, expiration)
}

func (r *RedisClient) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *redis.StringCmd {
	key := r.config.Prefix + source
	if r.config.Framework == "cluster" {
		return r.ClusterClient.BRPopLPush(ctx, key, destination, timeout)
	}
	return r.Client.BRPopLPush(ctx, key, destination, timeout)
}

func (r *RedisClient) LRem(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.LRem(ctx, key, count, value)
	}
	return r.Client.LRem(ctx, key, count, value)
}

func (r *RedisClient) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.ZRem(ctx, key, members)
	}
	return r.Client.ZRem(ctx, key, members)
}

func (r *RedisClient) Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd {
	channel = r.config.Prefix + channel
	if r.config.Framework == "cluster" {
		return r.ClusterClient.Publish(ctx, channel, message)
	}
	return r.Client.Publish(ctx, channel, message)
}

func (r *RedisClient) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.LRange(ctx, key, start, stop)
	}
	return r.Client.LRange(ctx, key, start, stop)
}

func (r *RedisClient) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	if r.config.Framework == "cluster" {
		return r.ClusterClient.Subscribe(ctx, channels...)
	}
	return r.Client.Subscribe(ctx, channels...)
}

func (r *RedisClient) LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.LPush(ctx, key, values...)
	}
	return r.Client.LPush(ctx, key, values...)
}

func (r *RedisClient) LLen(ctx context.Context, key string) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.LLen(ctx, key)
	}
	return r.Client.LLen(ctx, key)
}

func (r *RedisClient) LPop(ctx context.Context, key string) *redis.StringCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.LPop(ctx, key)
	}
	return r.Client.LPop(ctx, key)
}

func (r *RedisClient) RPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd {
	key = r.config.Prefix + key
	if r.config.Framework == "cluster" {
		return r.ClusterClient.RPush(ctx, key, values...)
	}
	return r.Client.RPush(ctx, key, values...)
}
