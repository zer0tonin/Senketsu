package model

import (
	"context"
	"fmt"
)

type TagRepository struct {
	repository Repository
}

func NewTagRepository(repository Repository) *TagRepository {
	return &TagRepository{
		repository: repository,
	}
}

func (r *TagRepository) Get(ctx context.Context, id string) (tag *Tag, err error) {
	result, err := r.repository.Get(ctx, id)
	if err != nil {
		return
	}

	tag, ok := result.(*Tag)
	if ok {
		return
	}
	return nil, fmt.Errorf("Failed to cast %s to Tag", result)
}

func (r *TagRepository) GetMany(ctx context.Context, ids []string) (tags []*Tag, err error) {
	results, err := r.repository.GetMany(ctx, ids)
	if err != nil {
		return
	}

	for _, result := range results {
		tag, ok := result.(*Tag)
		if !ok {
			return nil, fmt.Errorf("Failed to cast %s to Tag", result)
		}
		tags = append(tags, tag)
	}
	return
}


func (r *TagRepository) List(ctx context.Context) (tags []*Tag, err error) {
	results, err := r.repository.List(ctx)
	if err != nil {
		return
	}

	for _, result := range results {
		tag, ok := result.(*Tag)
		if !ok {
			return nil, fmt.Errorf("Failed to cast %s to Tag", result)
		}
		tags = append(tags, tag)
	}
	return
}

func (r *TagRepository) Save(ctx context.Context, tag *Tag) (error) {
	return r.repository.Save(ctx, tag)
}
