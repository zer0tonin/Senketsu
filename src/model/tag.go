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

func (t *Tag) getImages() []string {
	if t.Images == nil {
		t.Images = make([]string, 1)
	}
	return t.Images
}

func (t *Tag) AddImage(ctx context.Context, image *Image) error {
	t.Images = append(t.getImages(), image.ID)
	return nil
}

func (t *Tag) GetImages(ctx context.Context, image *Image) ([]*Image, error) {
	return S.ImageRepository.GetMany(ctx, t.getImages())
}

func (t *Tag) GetRandomImage(ctx context.Context) (*Image, error) {
	images := t.getImages()
	if len(images) == 0 {
		return nil, fmt.Errorf("No images found for this tag")
	}
	pick := rand.Int() % len(images)
	return S.ImageRepository.Get(ctx, t.Images[pick])
}

func (t *Tag) Save(ctx context.Context) (tag *Tag, err error) {
	return S.TagRepository.Save(ctx, t)
}
