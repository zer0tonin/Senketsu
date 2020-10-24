package model

import (
	"fmt"
	"net/http"
)

type ImageRepository interface {
	Get(ctx context.Context, id string) (*Image, error)
	GetMany(ctx context.Context, ids []string) ([]*Image, error)
	List(ctx context.Context) ([]*Image, error)
	Save(ctx context.Context, image *Image) (*Image, error)
}

type Image struct {
	ID  string `json:"id"`
	Uploader string `json:"uploader"`
	Tags []string `json:"tags"`
}

func FromRequest(ctx context.Context, r *http.Request) []*Image {
	fileHeaders, err := s.RequestParser.GetFileHeaders(r)
	var images = make([]*Image, 0, len(fileHeaders))
	if err != nil {
		fmt.Println(err)
		return images
	}

	for _, fileHeader := range fileHeaders {
		err = s.FileStorage.Upload(ctx, fileHeader)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Put image %s to S3\n", fileHeader.Filename)
		}
		imageCount := s.ImageIndex.IncrImageCount(ctx)
		image := &Image{
			ID: imageCount
			URI: fileHeader.FileName,
		}

		err = s.ImageIndex.AddImage(ctx, fileHeader.Filename, []string{"mytag"})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Put image %s to redis\n", fileHeader.Filename)
		}
	}
}
