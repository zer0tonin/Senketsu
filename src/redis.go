package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type RedisImageIndex struct {
	client *redis.Client
	prefix string
}

type Image struct {
	URI  string   `json:"uri"`
	Tags []string `json:"tags"`
}

func NewRedisImageIndex(redisAddress string, prefix string) *RedisImageIndex {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})
	return &RedisImageIndex{
		client: client,
		prefix: prefix,
	}
}

func (r *RedisImageIndex) AddImage(ctx context.Context, uri string, tags []string) (err error) {
	image := &Image{
		URI:  uri,
		Tags: tags,
	}
	ser, err := json.Marshal(image)
	if err != nil {
		return
	}

	incrCmd := r.client.Incr(ctx, fmt.Sprintf("%s:image_count", r.prefix))
	key, err := incrCmd.Result()
	if err != nil {
		return
	}
	status := r.client.Set(
		ctx,
		fmt.Sprintf("%s:image:%d", r.prefix, key),
		ser,
		0,
	)
	return status.Err()
}

func (r *RedisImageIndex) GetRandomImage(filePath string, tags []string) {
}
