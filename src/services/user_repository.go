package services

import (
	"context"
	"encoding/json"

	"github.com/zer0tonin/senketsu/src/model"
)

type UserRepository struct {
	database Database
}

func NewUserRepository(database Database) *UserRepository {
	return &UserRepository{
		database: database,
	}
}

func (r *UserRepository) unmarshal(ser []byte) (user *model.User, err error) {
	err = json.Unmarshal(ser, &user)
	return
}

func (r *UserRepository) marshal(user *model.User) (ser []byte, err error) {
	return json.Marshal(user)
}

func (r *UserRepository) Get(ctx context.Context, id string) (user *model.User, err error) {
	result, err := r.database.Get(ctx, id)
	if err != nil {
		return
	}
	return r.unmarshal(result)
}

func (r *UserRepository) GetMany(ctx context.Context, ids []string) (users []*model.User, err error) {
	results, err := r.database.GetMany(ctx, ids)
	if err != nil {
		return
	}

	for _, result := range results {
		if user, err := r.unmarshal(result); err != nil {
			return nil, err
		} else {
			users = append(users, user)
		}
	}
	return
}

func (r *UserRepository) List(ctx context.Context) (users []*model.User, err error) {
	results, err := r.database.List(ctx)
	if err != nil {
		return
	}

	for _, result := range results {
		if user, err := r.unmarshal(result); err != nil {
			return nil, err
		} else {
			users = append(users, user)
		}
	}
	return
}

func (r *UserRepository) Save(ctx context.Context, user *model.User) (error) {
	if ser, err := r.marshal(user); err != nil {
		return err
	} else {
		return r.database.Save(ctx, user.GetID(), ser)
	}
}
