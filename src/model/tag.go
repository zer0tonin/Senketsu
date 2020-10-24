package model

import (
	"context"
	"fmt"
	"math/rand"
)

type TagRepository interface {
	Add(ctx context.Context, tag *Tag) (*Tag, error)
	Get(ctx context.Context, id string) (*Tag, error)
	GetImages(ctx context.Context, tag *Tag) ([]string, error)
	List(ctx context.Context) ([]*Tag, error)
	Save(ctx context.Context, tag *Tag) (*Tag, error)
}

type Tag struct {
	Name   string `json:"name"`
	Images []string
}

func (t *Tag) GetImages(ctx context.Context) ([]string, error) {
	if t.Images == nil {
		images, err := S.TagRepository.GetImages(ctx, t)
		if err != nil {
			return nil, err
		}
		t.Images = images
	}
	return t.Images, nil
}

func (t *Tag) AddImage(ctx context.Context, image *Image) error {
	images, err := S.TagRepository.GetImages(ctx, t)
	if err != nil {
		return err
	}
	t.Images = append(images, image.ID)
	return nil
}

func (t *Tag) GetRandomImage(ctx context.Context) (*Image, error) {
	images, err := t.GetImages(ctx)
	if err != nil {
		return nil, err
	}
	if len(images) == 0 {
		return nil, fmt.Errorf("No images found for this tag")
	}
	pick := rand.Int() % len(images)
	return S.ImageRepository.Get(ctx, t.Images[pick])
}

func (t *Tag) Save(ctx context.Context) (*Tag, error) {
	tag, err := S.TagRepository.Save(ctx, t)
	// save images
	return tag, err
}
