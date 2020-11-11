package model

import (
	"context"
	"html/template"
	"mime/multipart"
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
	ParseForm(r *http.Request) (result []*multipart.FileHeader, err error)
}

type FileStorage interface {
	WriteFile(ctx context.Context, fileHeader *multipart.FileHeader) error
}

type Services struct {
	ImageRepository ImageRepository
	TagRepository   TagRepository
	RequestParser   RequestParser
	FileStorage     FileStorage
	Templates       map[string]*template.Template
}

var S Services
