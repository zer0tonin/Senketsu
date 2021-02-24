package model

import (
	"context"
	"net/http"
	"net/url"
)

type Entity interface {
	GetID() string
}

type ImageRepository interface {
	Get(ctx context.Context, id string) (*Image, error)
	GetMany(ctx context.Context, ids []string) ([]*Image, error)
	List(ctx context.Context) ([]*Image, error)
	Save(ctx context.Context, image *Image) (error)
}

type TagRepository interface {
	Get(ctx context.Context, id string) (*Tag, error)
	GetMany(ctx context.Context, ids []string) ([]*Tag, error)
	List(ctx context.Context) ([]*Tag, error)
	Save(ctx context.Context, tag *Tag) (error)
}

type UserRepository interface {
	Get(ctx context.Context, id string) (*User, error)
	GetMany(ctx context.Context, ids []string) ([]*User, error)
	List(ctx context.Context) ([]*User, error)
	Save(ctx context.Context, user *User) (error)
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
	ImageRepository ImageRepository
	TagRepository TagRepository
	UserRepository UserRepository
	RequestParser   RequestParser
	FileStorage     FileStorage
	Views           Views
	AuthenticationProvider AuthenticationProvider
}

var S Services
