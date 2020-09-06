package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

type Services struct {
	MultipartParser *MultipartParser
	S3Client        *S3Client
	Templates       map[string]*template.Template
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

	minioClient, err := minio.New(
		viper.GetString("s3.endpoint"),
		&minio.Options{
			Creds: credentials.NewStaticV4(
				viper.GetString("s3.accessKeyID"),
				viper.GetString("s3.secretAccessKey"),
				"",
			),
			Secure: viper.GetBool("s3.useSSL"),
		},
	)
	if err != nil {
		panic("Failed to create S3 client")
	}

	multipartParser := NewMultipartParser(viper.GetInt64("multipartParser.maxMem"))
	s3Client := NewS3Client(
		viper.GetString("s3.bucket"),
		minioClient,
	)

	s = Services{
		MultipartParser: multipartParser,
		S3Client:        s3Client,
		Templates:       templates,
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

		fileHeaders, err := s.MultipartParser.GetFileHeaders(r)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, fileHeader := range fileHeaders {
			s.S3Client.Upload(ctx, fileHeader)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Put Object %s to S3\n", fileHeader.Filename)
			}
		}
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
