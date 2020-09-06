package main

import (
	"fmt"
	"html/template"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
)

type Services struct {
	S3Client  *minio.Client
	Templates map[string]*template.Template
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

	templates := make(map[string]*template.Template)
	templates["index"] = template.Must(template.ParseFiles("./templates/index.html"))

	s = Services{
		S3Client:  minioClient,
		Templates: templates,
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

		_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Parsed MediaType")

		mr := multipart.NewReader(r.Body, params["boundary"])
		form, err := mr.ReadForm(134217728) // 128MiB
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Parsed Multipart Form")

		for _, fileHeaders := range form.File {
			for _, fileHeader := range fileHeaders {
				file, err := fileHeader.Open()
				if err != nil {
					fmt.Println(err)
					continue
				}
				_, err = s.S3Client.PutObject(
					ctx,
					"senketsu-test",
					fileHeader.Filename,
					file,
					fileHeader.Size,
					minio.PutObjectOptions{},
				)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Printf("Put Object %s to S3\n", fileHeader.Filename)
				}
			}
		}
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
