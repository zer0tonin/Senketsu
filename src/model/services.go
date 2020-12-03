package model

import (
	"context"
	"net/http"
	"net/url"
)

type Entity interface {
	GetID() string
}

type Repository interface {
	Get(ctx context.Context, id string) (Entity, error)
	GetMany(ctx context.Context, ids []string) ([]Entity, error)
	List(ctx context.Context) ([]Entity, error)
	Save(ctx context.Context, entity Entity) (error)
}

type RequestParser interface {
	ParseForm(r *http.Request) (result []*Image, errs []error)
}

type FileStorage interface {
	WriteFile(ctx context.Context, image *Image) error
	GetURL(image *Image) (*url.URL, error)
}

type Views interface {
	Index(w http.ResponseWriter)
	UploadResults(w http.ResponseWriter, images []*Image)
}

type AuthenticationToken interface {
	TokenResponse(w http.ResponseWriter)
}

type AuthenticationProvider interface {
	Request(ctx context.Context, u *User, w http.ResponseWriter)
	Callback(ctx context.Context, u *User, w http.ResponseWriter) AuthenticationToken
}

type Services struct {
	*ImageRepository
	*TagRepository
	*UserRepository
	RequestParser   RequestParser
	FileStorage     FileStorage
	Views           Views
	AuthenticationProvider AuthenticationProvider
}

var S Services
