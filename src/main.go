package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"

	"github.com/zer0tonin/senketsu/src/model"
	"github.com/zer0tonin/senketsu/src/services"
	"github.com/zer0tonin/senketsu/src/views"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to load config")
	}

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

	redisTagRepository := services.NewRedisTagRepository(
		redisClient,
		viper.GetString("redis.prefix"),
	)

	model.S = model.Services{
		RequestParser:   multipartParser,
		ImageRepository: redisImageRepository,
		TagRepository:   redisTagRepository,
		FileStorage:     s3Storage,
	}

	tags := viper.GetStringSlice("defaultTags")
	for _, tagName := range tags {
		tag := &model.Tag{
			Name: tagName,
		}
		_, err := tag.Save(context.Background())
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	http.Handle("/", views.BaseHandler())
	http.Handle("/images/", views.ImagesHandler())
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
