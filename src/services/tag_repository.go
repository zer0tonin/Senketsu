package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/gosimple/slug"

	"github.com/zer0tonin/senketsu/src/model"
)

type RedisTagRepository struct {
	client *redis.Client
	prefix string
}

func (r *RedisTagRepository) makeKey() string {
	return fmt.Sprintf("%s:tags", r.prefix)
}

func (r *RedisTagRepository) makeImageKey(id string) string {
	return fmt.Sprintf("%s:tag:%d:images", r.prefix, slug.Make(id))
}

func (r *RedisTagRepository) Add(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	return r.Save(ctx, tag)
}

func (r *RedisTagRepository) Get(ctx context.Context, id string) (tag *model.Tag, err error) {
	result := r.client.HGet(ctx, r.makeKey(), id)
	err = result.Err()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(result.Val()), tag)
	return tag, err
}

func (r *RedisTagRepository) GetImages(ctx context.Context, tag *model.Tag) (images []string, err error) {
	result := r.client.Get(ctx, r.makeImageKey(tag.Name))
	err = result.Err()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(result.Val()), images)
	return
}

func (r *RedisTagRepository) List(ctx context.Context) (tags []*model.Tag, err error) {
	result := r.client.HGetAll(ctx, r.makeKey())
	err = result.Err()
	if err != nil {
		return
	}

	for _, ser := range result.Val() {
		var tag *model.Tag
		err = json.Unmarshal([]byte(ser), tag)
	}
	return
}

func (r *RedisTagRepository) Save(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	ser, err := json.Marshal(tag)
	if err != nil {
		return nil, err
	}

	result := r.client.HSet(
		ctx,
		r.makeKey(),
		slug.Make(tag.Name),
		ser,
		0,
	)
	return tag, result.Err()
}
