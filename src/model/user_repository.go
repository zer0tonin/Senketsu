package model

import (
	"context"
	"fmt"
)

type UserRepository struct {
	repository Repository
}

func NewUserRepository(repository Repository) *UserRepository {
	return &UserRepository{
		repository: repository,
	}
}

func (r *UserRepository) Get(ctx context.Context, id string) (user *User, err error) {
	result, err := r.repository.Get(ctx, id)
	if err != nil {
		return
	}

	user, ok := result.(*User)
	if ok {
		return
	}
	return nil, fmt.Errorf("Failed to cast %s to User", result)
}

func (r *UserRepository) GetMany(ctx context.Context, ids []string) (users []*User, err error) {
	results, err := r.repository.GetMany(ctx, ids)
	if err != nil {
		return
	}

	for _, result := range results {
		user, ok := result.(*User)
		if !ok {
			return nil, fmt.Errorf("Failed to cast %s to User", result)
		}
		users = append(users, user)
	}
	return
}

func (r *UserRepository) List(ctx context.Context) (users []*User, err error) {
	results, err := r.repository.List(ctx)
	if err != nil {
		return
	}

	for _, result := range results {
		user, ok := result.(*User)
		if !ok {
			return nil, fmt.Errorf("Failed to cast %s to User", result)
		}
		users = append(users, user)
	}
	return
}

func (r *UserRepository) Save(ctx context.Context, user *User) (error) {
	return r.repository.Save(ctx, user)
}
