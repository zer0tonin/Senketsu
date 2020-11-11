package services

import (
	"context"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Storage struct {
	bucketName string
	client     *minio.Client
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

	return &S3Storage{
		bucketName: bucketName,
		client:     minioClient,
	}
}

func (s *S3Storage) WriteFile(ctx context.Context, fileHeader *multipart.FileHeader) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	_, err = s.client.PutObject(
		ctx,
		s.bucketName,
		fileHeader.Filename,
		file,
		fileHeader.Size,
		minio.PutObjectOptions{},
	)
	return err
}
