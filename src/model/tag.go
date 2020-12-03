package model

import (
	"context"
	"fmt"
	"math/rand"
)

type Tag struct {
	Name   string   `json:"name"`
	Images []string `json:"images"`
}

func NewTag(name string) *Tag {
	return &Tag{
		Name: name,
		Images: make([]string, 0),
	}
}

func (t *Tag) GetID() string {
	return t.Name
}

func (t *Tag) AddImage(image *Image) {
	// TODO : could be optimized by making t.Images a set
	for _, i := range t.Images {
		if i == image.ID {
			return
		}
	}
	t.Images = append(t.Images, image.ID)
	return
}

func (t *Tag) GetImages(ctx context.Context, image *Image) ([]*Image, error) {
	return S.ImageRepository.GetMany(ctx, t.Images)
}

func (t *Tag) GetRandomImage(ctx context.Context) (*Image, error) {
	images := t.Images
	if len(images) == 0 {
		return nil, fmt.Errorf("No images found for this tag")
	}
	pick := rand.Int() % len(images)
	return S.ImageRepository.Get(ctx, t.Images[pick])
}

func (t *Tag) Save(ctx context.Context) (err error) {
	return S.TagRepository.Save(ctx, t)
}
