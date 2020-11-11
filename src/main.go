package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"

	"github.com/zer0tonin/senketsu/src/model"
	"github.com/zer0tonin/senketsu/src/services"
)

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

	multipartParser := services.NewMultipartParser(viper.GetInt64("multipartParser.maxMem"))

	s3Storage := services.NewS3Storage(
		viper.GetString("s3.endpoint"),
		viper.GetString("s3.accessKeyID"),
		viper.GetString("s3.secretAccessKey"),
		viper.GetBool("s3.useSSL"),
		viper.GetString("s3.bucket"),
	)

	redisClient := redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.address"),
	})

	redisImageRepository := services.NewRedisImageRepository(
		redisClient,
		viper.GetString("redis.prefix"),
	)

	model.S = model.Services{
		RequestParser:   multipartParser,
		ImageRepository: redisImageRepository,
		FileStorage:     s3Storage,
		Templates:       templates,
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		model.S.Templates["index"].Execute(w, nil)
	})

	http.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/index.css")
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		model.NewImageFromRequest(ctx, r)
		// TODO: don't ignore errors
		// TODO: use different error types for the bad request, internal error, etc
		// TODO: redirect to image view
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
