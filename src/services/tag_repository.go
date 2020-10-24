package services

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
	return fmt.Sprintf("%s:tag:%s:images", r.prefix, slug.Make(id))
}

func (r *RedisTagRepository) Get(ctx context.Context, id string) (tag *model.Tag, err error) {
	result := r.client.HGet(ctx, r.makeKey(), id)
	if err = result.Err(); err != nil {
		return
	}

	err = json.Unmarshal([]byte(result.Val()), tag)
	return
}

func (r *RedisTagRepository) List(ctx context.Context) (tags []*model.Tag, err error) {
	result := r.client.HGetAll(ctx, r.makeKey())
	if err = result.Err(); err != nil {
		return
	}

	for _, ser := range result.Val() {
		var tag *model.Tag
		if err := json.Unmarshal([]byte(ser), tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return
}

func (r *RedisTagRepository) Save(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	if ser, err := json.Marshal(tag); err != nil {
		return nil, err
	} else {
		result := r.client.HSet(
			ctx,
			r.makeKey(),
			slug.Make(tag.Name),
			ser,
			0,
		)
		return tag, result.Err()
	}
}
