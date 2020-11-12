package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/gosimple/slug"

	"github.com/zer0tonin/senketsu/src/model"
)

type RedisImageRepository struct {
	client *redis.Client
	prefix string
}

func NewRedisImageRepository(client *redis.Client, prefix string) *RedisImageRepository {
	return &RedisImageRepository{
		client: client,
		prefix: prefix,
	}
}

func (r *RedisImageRepository) makeKey() string {
	return fmt.Sprintf("%s:images", r.prefix)
}

func (r *RedisImageRepository) makeTagsKey(id string) string {
	return fmt.Sprintf("%s:images:%s:tags", r.prefix, slug.Make(id))
}

func (r *RedisImageRepository) Get(ctx context.Context, id string) (image *model.Image, err error) {
	result := r.client.HGet(ctx, r.makeKey(), id)
	if err = result.Err(); err != nil {
		return
	}

	image = &model.Image{}
	err = json.Unmarshal([]byte(result.Val()), image)
	return
}

func (r *RedisImageRepository) GetMany(ctx context.Context, ids []string) (images []*model.Image, err error) {
	result := r.client.HMGet(ctx, r.makeKey(), ids...)
	if err = result.Err(); err != nil {
		return
	}

	for _, ser := range result.Val() {
		var image *model.Image
		if err := json.Unmarshal([]byte(ser.(string)), image); err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return
}

func (r *RedisImageRepository) List(ctx context.Context) (images []*model.Image, err error) {
	result := r.client.HGetAll(ctx, r.makeKey())
	if err = result.Err(); err != nil {
		return
	}

	for _, ser := range result.Val() {
		image := &model.Image{}
		if err := json.Unmarshal([]byte(ser), image); err != nil {
			return nil, err
		}
		images = append(images, image)
	}
	return
}

func (r *RedisImageRepository) Save(ctx context.Context, image *model.Image) (*model.Image, error) {
	if ser, err := json.Marshal(image); err != nil {
		return nil, err
	} else {
		result := r.client.HSet(
			ctx,
			r.makeKey(),
			image.ID,
			ser,
		)
		return image, result.Err()
	}
}
