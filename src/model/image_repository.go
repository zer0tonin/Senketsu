package model

import (
	"context"
	"fmt"
)

type ImageRepository struct {
	repository Repository
}

func NewImageRepository(repository Repository) *ImageRepository {
	return &ImageRepository{
		repository: repository,
	}
}

func (r *ImageRepository) Get(ctx context.Context, id string) (image *Image, err error) {
	result, err := r.repository.Get(ctx, id)
	if err != nil {
		return
	}

	image, ok := result.(*Image)
	if ok {
		return
	}
	return nil, fmt.Errorf("Failed to cast %s to Image", result)
}

func (r *ImageRepository) GetMany(ctx context.Context, ids []string) (images []*Image, err error) {
	results, err := r.repository.GetMany(ctx, ids)
	if err != nil {
		return
	}

	for _, result := range results {
		image, ok := result.(*Image)
		if !ok {
			return nil, fmt.Errorf("Failed to cast %s to Image", result)
		}
		images = append(images, image)
	}
	return
}

func (r *ImageRepository) List(ctx context.Context) (images []*Image, err error) {
	results, err := r.repository.List(ctx)
	if err != nil {
		return
	}

	for _, result := range results {
		image, ok := result.(*Image)
		if !ok {
			return nil, fmt.Errorf("Failed to cast %s to Image", result)
		}
		images = append(images, image)
	}
	return
}

func (r *ImageRepository) Save(ctx context.Context, image *Image) (error) {
	return r.repository.Save(ctx, image)
}
