package model

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type Image struct {
	ID        string `json:"id"`
	Extension string `json:"id"`
	//Uploader string   `json:"uploader"`
	//Tags     []string `json:"tags"`
	Reader io.Reader
	Size   int64
}

func NewImageFromRequest(ctx context.Context, r *http.Request) (images []*Image, errs []error) {
	images, errs = S.RequestParser.ParseForm(r)
	if len(errs) != 0 {
		return nil, errs
	}

	for _, image := range images {
		imageID, err := uuid.NewRandom()
		if err != nil {
			errs = append(errs, err)
		}
		image.ID = imageID.String()
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
		_, err = image.Save(ctx)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return images, errs
}

func (i *Image) Save(ctx context.Context) (*Image, error) {
	return S.ImageRepository.Save(ctx, i)
}

func (i *Image) GetURI() string {
	return S.FileStorage.GetURI(i)
}

func (i *Image) GetFilename() string {
	return fmt.Sprintf("%s.%s", i.ID, i.Extension)
}
