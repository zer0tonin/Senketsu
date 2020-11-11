package model

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Image struct {
	ID string `json:"id"`
	//Uploader string   `json:"uploader"`
	//Tags     []string `json:"tags"`
}

func NewImageFromRequest(ctx context.Context, r *http.Request) (*Image, error) {
	fileHeaders, err := S.RequestParser.ParseForm(r)
	if err != nil {
		return nil, err
	}

	if len(fileHeaders) != 1 {
		return nil, fmt.Errorf("Please upload a single image")
	}

	imageID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	image := &Image{
		ID: imageID.String(),
	}

	fileHeader := fileHeaders[0]
	err = S.FileStorage.WriteFile(ctx, fileHeader)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Put image %s to S3\n", fileHeader.Filename)
	}

	image, err = image.Save(ctx)
	if err != nil {
		return nil, err
	}
	return image, nil
}

func (i *Image) Save(ctx context.Context) (*Image, error) {
	return S.ImageRepository.Save(ctx, i)
}
