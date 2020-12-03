package services

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"

	"github.com/zer0tonin/senketsu/src/model"
)

type RedisRepository struct {
	client *redis.Client
	prefix string
}

func NewRedisRepository(client *redis.Client, prefix string) *RedisRepository {
	return &RedisRepository{
		client: client,
		prefix: prefix,
	}
}

func (r *RedisRepository) Get(ctx context.Context, id string) (value model.Entity, err error) {
	result := r.client.HGet(ctx, r.prefix, id)
	if err = result.Err(); err != nil {
		return
	}

	err = json.Unmarshal([]byte(result.Val()), &value)
	return
}

func (r *RedisRepository) GetMany(ctx context.Context, ids []string) (values []model.Entity, err error) {
	result := r.client.HMGet(ctx, r.prefix, ids...)
	if err = result.Err(); err != nil {
		return
	}

	for _, ser := range result.Val() {
		var value model.Entity
		if err := json.Unmarshal([]byte(ser.(string)), &value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return
}

func (r *RedisRepository) List(ctx context.Context) (values []model.Entity, err error) {
	result := r.client.HGetAll(ctx, r.prefix)
	if err = result.Err(); err != nil {
		return
	}

	for _, ser := range result.Val() {
		var value model.Entity
		if err := json.Unmarshal([]byte(ser), value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return
}

func (r *RedisRepository) Save(ctx context.Context, value model.Entity) error {
	ser, err := json.Marshal(value)
	if err != nil {
		return err
	}
	result := r.client.HSet(
		ctx,
		r.prefix,
		value.GetID(),
		ser,
	)
	return result.Err()
}
