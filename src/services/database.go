package services

import (
	"context"
)

type Database interface {
	Get(ctx context.Context, id string) ([]byte, error)
	GetMany(ctx context.Context, ids []string) ([][]byte, error)
	List(ctx context.Context) ([][]byte, error)
	Save(ctx context.Context, id string, value []byte) (error)
}
