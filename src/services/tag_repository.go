package services

import (
	"context"
	"encoding/json"

	"github.com/zer0tonin/senketsu/src/model"
)

type TagRepository struct {
	database Database
}

func NewTagRepository(database Database) *TagRepository {
	return &TagRepository{
		database: database,
	}
}

func (r *TagRepository) unmarshal(ser []byte) (tag *model.Tag, err error) {
	err = json.Unmarshal(ser, &tag)
	return nil, err
}

func (r *TagRepository) marshal(tag *model.Tag) (ser []byte, err error) {
	return json.Marshal(tag)
}

func (r *TagRepository) Get(ctx context.Context, id string) (tag *model.Tag, err error) {
	result, err := r.database.Get(ctx, id)
	if err != nil {
		return
	}
	return r.unmarshal(result)
}

func (r *TagRepository) GetMany(ctx context.Context, ids []string) (tags []*model.Tag, err error) {
	results, err := r.database.GetMany(ctx, ids)
	if err != nil {
		return
	}

	for _, result := range results {
		if tag, err := r.unmarshal(result); err != nil {
			return nil, err
		} else {
			tags = append(tags, tag)
		}
	}
	return
}


func (r *TagRepository) List(ctx context.Context) (tags []*model.Tag, err error) {
	results, err := r.database.List(ctx)
	if err != nil {
		return
	}

	for _, result := range results {
		if tag, err := r.unmarshal(result); err != nil {
			return nil, err
		} else {
			tags = append(tags, tag)
		}
	}
	return
}

func (r *TagRepository) Save(ctx context.Context, tag *model.Tag) (error) {
	if ser, err := r.marshal(tag); err != nil {
		return err
	} else {
		return r.database.Save(ctx, tag.GetID(), ser)
	}
}
