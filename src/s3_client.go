package main

import (
	"context"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Client struct {
	bucketName string
	client     *minio.Client
}

func NewS3Client(
	endpoint string,
	accessKeyID string,
	secretAccessKey string,
	useSSL bool,
	bucketName string,
) *S3Client {
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

	return &S3Client{
		bucketName: bucketName,
		client:     minioClient,
	}
}

func (s *S3Client) Upload(ctx context.Context, fileHeader *multipart.FileHeader) error {
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
