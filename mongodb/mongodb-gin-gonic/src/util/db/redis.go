package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisUtil struct {
	addr     string
	password string
	db       int
	client   *redis.Client
}

func NewRedisConnection(addr, password string, db int) *RedisUtil {
	return &RedisUtil{
		addr:     addr,
		password: password,
		db:       db,
	}
}

func (o *RedisUtil) connect() {
	if o.client == nil {
		o.client = redis.NewClient(
			&redis.Options{
				Addr:     o.addr,
				Password: o.password,
				DB:       o.db,
			})
	}
}

func (o *RedisUtil) close() error {
	return o.client.Close()
}

func (o *RedisUtil) Get(key string) (string, error) {
	o.connect()
	defer o.close()

	return o.client.Get(context.Background(), key).Result()
}

func (o *RedisUtil) Set(key, value string, expiration time.Duration) error {
	o.connect()
	defer o.close()

	return o.client.Set(context.Background(), key, value, expiration).Err()
}
