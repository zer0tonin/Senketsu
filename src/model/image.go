package model

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type Image struct {
	ID        string `json:"id"`
	Extension string `json:"extension"`
	//Uploader string   `json:"uploader"`
	Tags     []string `json:"tags"`
	Reader io.Reader `json:"-"`
	Size   int64     `json:"size"`
}

func NewImageFromRequest(ctx context.Context, r *http.Request) (images []*Image, errs []error) {
	images, errs = S.RequestParser.ParseForm(r)
	if len(errs) != 0 {
		return nil, errs
	}

	if len(errs) != 0 {
		return nil, errs
	}

	for _, image := range images {
		err := S.FileStorage.WriteFile(ctx, image)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		fmt.Printf("Put image %s to S3\n", image.GetFilename())
		// TODO: risks of orphan images
		err = image.Save(ctx)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return images, errs
}

func (i *Image) GetID() string {
	if i.ID == "" {
		imageID, err := uuid.NewRandom()
		if err != nil {
			panic("UUID generation broken")
		}
		i.ID = imageID.String()
	}
	return i.ID
}

func (i *Image) Save(ctx context.Context) (error) {
	for _, tagName := range i.Tags {
		tag, _:= S.TagRepository.Get(ctx, tagName)
		tag.AddImage(i)
		tag.Save(ctx)
	}
	return S.ImageRepository.Save(ctx, i)
}

func (i *Image) GetStorageURL() (*url.URL, error) {
	return S.FileStorage.GetURL(i)
}

func (i *Image) GetPublicURL() string {
	host := viper.GetString("host")
	if host[len(host)-1] == '/' {
		return fmt.Sprintf("%sfiles/%s.%s", host, i.ID, i.Extension)
	} else {
		return fmt.Sprintf("%s/files/%s.%s", host, i.ID, i.Extension)
	}
}

func (i *Image) GetFilename() string {
	return fmt.Sprintf("%s.%s", i.GetID(), i.Extension)
}
