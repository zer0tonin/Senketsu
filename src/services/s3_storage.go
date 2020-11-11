package services

import (
	"context"
	"fmt"
	"net/url"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/zer0tonin/senketsu/src/model"
)

type S3Storage struct {
	bucketName string
	client     *minio.Client
	baseURI    string
}

func NewS3Storage(
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
	useSSL bool,
	bucketName string,
) *S3Storage {
	minioClient, err := minio.New(
		endpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				accessKeyID,
				secretAccessKey,
				"",
			),
			Secure: useSSL,
		},
	)
	if err != nil {
		panic("Failed to create S3 client")
	}

	var prefix string
	if useSSL {
		prefix = "https"
	} else {
		prefix = "http"
	}

	return &S3Storage{
		bucketName: bucketName,
		client:     minioClient,
		baseURI:    fmt.Sprintf("%s://%s/%s", prefix, endpoint, bucketName),
	}
}

func (s *S3Storage) WriteFile(
	ctx context.Context,
	image *model.Image,
) error {
	_, err := s.client.PutObject(
		ctx,
		s.bucketName,
		image.GetFilename(),
		image.Reader,
		image.Size, //fileHeader.Size,
		minio.PutObjectOptions{},
	)
	return err
}

func (s *S3Storage) GetURL(image *model.Image) (*url.URL, error) {
	return url.Parse(fmt.Sprintf("%s/%s", s.baseURI, image.GetFilename()))
}
