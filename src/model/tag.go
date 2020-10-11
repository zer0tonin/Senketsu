package model

import (
	"context"
	"math/rand"
)

type TagRepository interface {
	Add(ctx context.Context, tag *Tag) (*Tag, error)
	Save(ctx context.Context, tag *Tag) (*Tag, error)
	GetImageIDs(ctx context.Context, tag *Tag) ([]string, error)
}

type Tag struct {
	ID              *int   `json:"id"`
	Name            string `json:"name"`
	tagRepository   TagRepository
	imageRepository ImageRepository
}

func NewTag(
	ctx context.Context,
	tagRepository TagRepository,
	imageRepository ImageRepository,
	name string,
) (*Tag, error) {
	tag := &Tag{
		Name:          name,
		tagRepository: tagRepository,
	}
	return tag.Save(ctx)
}

func (t *Tag) Save(ctx context.Context) (*Tag, error) {
	if t.ID == nil {
		return t.tagRepository.Add(ctx, t)
	}
	return t.tagRepository.Save(ctx, t)
}

func (t *Tag) GetRandomImage(ctx context.Context) (*Image, error) {
	imagesIDs, err := t.tagRepository.GetImageIDs(ctx, t)
	if err != nil {
		return nil, err
	}
	pick := rand.Int() % len(imagesIDs)
	return t.imageRepository.Get(ctx, imagesIDs[pick])
}
