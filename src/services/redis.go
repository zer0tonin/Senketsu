package services

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
	prefix string
}

func NewRedis(client *redis.Client, prefix string) *Redis {
	return &Redis{
		client: client,
		prefix: prefix,
	}
}

func (r *Redis) Get(ctx context.Context, id string) (value []byte, err error) {
	result := r.client.HGet(ctx, r.prefix, id)
	if err = result.Err(); err != nil {
		return
	}

	return []byte(result.Val()), nil
}

func (r *Redis) GetMany(ctx context.Context, ids []string) (values [][]byte, err error) {
	result := r.client.HMGet(ctx, r.prefix, ids...)
	if err = result.Err(); err != nil {
		return
	}

	for _, ser := range result.Val() {
		values = append(values, []byte(ser.(string)))
	}
	return
}

func (r *Redis) List(ctx context.Context) (values [][]byte, err error) {
	result := r.client.HGetAll(ctx, r.prefix)
	if err = result.Err(); err != nil {
		return
	}

	for _, ser := range result.Val() {
		values = append(values, []byte(ser))
	}
	return
}

func (r *Redis) Save(ctx context.Context, id string, value []byte) error {
	result := r.client.HSet(
		ctx,
		r.prefix,
		id,
		value,
	)
	return result.Err()
}
