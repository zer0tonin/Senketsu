package main

import (
	"context"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/http"

	"github.com/spf13/viper"
)

type RequestParser interface {
	GetFileHeaders(r *http.Request) (result []*multipart.FileHeader, err error)
}

type ImageIndex interface {
	AddImage(ctx context.Context, uri string, tags []string) (err error)
}

type FileStorage interface {
	Upload(ctx context.Context, fileHeader *multipart.FileHeader) (err error)
}

type Services struct {
	ImageIndex    ImageIndex
	FileStorage   FileStorage
	RequestParser RequestParser
	Templates     map[string]*template.Template
}

var s Services

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to load config")
	}

	templates := make(map[string]*template.Template)
	templates["index"] = template.Must(template.ParseFiles("./templates/index.html"))

	multipartParser := NewMultipartParser(viper.GetInt64("multipartParser.maxMem"))

	s3Client := NewS3Client(
		viper.GetString("s3.endpoint"),
		viper.GetString("s3.accessKeyID"),
		viper.GetString("s3.secretAccessKey"),
		viper.GetBool("s3.useSSL"),
		viper.GetString("s3.bucket"),
	)

	redisImageIndex := NewRedisImageIndex(
		viper.GetString("redis.address"),
		viper.GetString("redis.prefix"),
	)

	s = Services{
		RequestParser: multipartParser,
		FileStorage:   s3Client,
		Templates:     templates,
		ImageIndex:    redisImageIndex,
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.Templates["index"].Execute(w, nil)
	})

	http.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/index.css")
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		fileHeaders, err := s.RequestParser.GetFileHeaders(r)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, fileHeader := range fileHeaders {
			err = s.FileStorage.Upload(ctx, fileHeader)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Put image %s to S3\n", fileHeader.Filename)
			}
			err = s.ImageIndex.AddImage(ctx, fileHeader.Filename, []string{"mytag"})
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Put image %s to redis\n", fileHeader.Filename)
			}
		}
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
