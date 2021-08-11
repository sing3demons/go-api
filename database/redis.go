package database

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sing3demons/api/models"
)

type redisCache struct {
	host    string
	db      int
	expires time.Duration
}

type RedisCache interface {
	Set(key string, value interface{}) error
	Get(key string) ([]models.Product, error)
}

func NewRedisCache(host string, db int, exp time.Duration) RedisCache {
	return &redisCache{host: host,
		db:      db,
		expires: exp}
}

func (cache *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cache.host,
		Password: "",
		DB:       cache.db,
	})
}

func (cache *redisCache) Set(key string, value interface{}) error {
	rdb := cache.getClient()
	ctx := context.Background()

	json, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	err = rdb.Set(ctx, key, json, cache.expires*time.Second).Err()
	if err != nil {
		fmt.Printf("set error: %v", err)
		return err
	}

	return nil
}
func (cache *redisCache) Get(key string) ([]models.Product, error) {
	rdb := cache.getClient()
	ctx := context.Background()
	val, err := rdb.Get(ctx, key).Result()
	
	if err != nil {
		return nil, err
	}

	product := []models.Product{}
	if err := json.Unmarshal([]byte(val), &product); err != nil {
		panic(err)
	}

	return product, nil
}
