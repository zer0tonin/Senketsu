package model

import (
	"context"
	"html/template"
	"net/http"
)

type TagRepository interface {
	Get(ctx context.Context, id string) (*Tag, error)
	List(ctx context.Context) ([]*Tag, error)
	Save(ctx context.Context, tag *Tag) (*Tag, error)
}

type ImageRepository interface {
	Get(ctx context.Context, id string) (*Image, error)
	GetMany(ctx context.Context, ids []string) ([]*Image, error)
	List(ctx context.Context) ([]*Image, error)
	Save(ctx context.Context, image *Image) (*Image, error)
}

type RequestParser interface {
	ParseForm(r *http.Request) (result []*Image, errs []error)
}

type FileStorage interface {
	WriteFile(ctx context.Context, image *Image) error
	GetURI(image *Image) string
}

type Services struct {
	ImageRepository ImageRepository
	TagRepository   TagRepository
	RequestParser   RequestParser
	FileStorage     FileStorage
	Templates       map[string]*template.Template
}

var S Services
