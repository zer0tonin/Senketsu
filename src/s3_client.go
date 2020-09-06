package main

import (
	"context"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type S3Client struct {
	bucketName string
	client     *minio.Client
}

func NewS3Client(bucketName string, client *minio.Client) *S3Client {
	return &S3Client{
		bucketName: bucketName,
		client:     client,
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
