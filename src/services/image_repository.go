package services

import (
	"context"
	"encoding/json"

	"github.com/zer0tonin/senketsu/src/model"
)

type ImageRepository struct {
	database Database
}

func NewImageRepository(database Database) *ImageRepository {
	return &ImageRepository{
		database: database,
	}
}

func (r *ImageRepository) unmarshal(ser []byte) (image *model.Image, err error) {
	err = json.Unmarshal(ser, &image)
	return
}

func (r *ImageRepository) marshal(image *model.Image) (ser []byte, err error) {
	return json.Marshal(image)
}

func (r *ImageRepository) Get(ctx context.Context, id string) (image *model.Image, err error) {
	result, err := r.database.Get(ctx, id)
	if err != nil {
		return
	}
	return r.unmarshal(result)
}

func (r *ImageRepository) GetMany(ctx context.Context, ids []string) (images []*model.Image, err error) {
	results, err := r.database.GetMany(ctx, ids)
	if err != nil {
		return
	}

	for _, result := range results {
		if image, err := r.unmarshal(result); err != nil {
			return nil, err
		} else {
			images = append(images, image)
		}
	}
	return
}

func (r *ImageRepository) List(ctx context.Context) (images []*model.Image, err error) {
	results, err := r.database.List(ctx)
	if err != nil {
		return
	}

	for _, result := range results {
		if image, err := r.unmarshal(result); err != nil {
			return nil, err
		} else {
			images = append(images, image)
		}
	}
	return
}

func (r *ImageRepository) Save(ctx context.Context, image *model.Image) (error) {
	if ser, err := r.marshal(image); err != nil {
		return err
	} else {
		return r.database.Save(ctx, image.GetID(), ser)
	}
}
