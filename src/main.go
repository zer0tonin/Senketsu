package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httputil"
	"regexp"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
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
	}
}

func main() {
	templates := make(map[string]*template.Template)
	templates["index"] = template.Must(template.ParseFiles("./templates/index.html"))
	templates["upload"] = template.Must(template.ParseFiles("./templates/upload.html"))

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := templates["index"].Execute(w, nil)
		if err != nil {
			fmt.Println(err)
		}
	})

	r.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/index.css")
	})

	r.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		images, errs := model.NewImageFromRequest(r.Context(), r)
		err := templates["upload"].Execute(
			w,
			map[string]interface{}{
				"images": images,
				"errs":   errs,
			},
		)
		if err != nil {
			fmt.Println(err)
		}
	})

	imagesRegexp := regexp.MustCompile(`.+\/(.+)\.gif`)
	r.HandleFunc("/files/{path}", func(w http.ResponseWriter, r *http.Request) {
		matches := imagesRegexp.FindStringSubmatch(r.URL.Path)
		if len(matches) == 2 {
			image, err := model.S.ImageRepository.Get(r.Context(), matches[1])
			if err != nil {
				fmt.Println(err)
				return
			}
			if image == nil {
				fmt.Println("404") //TODO
				return
			}

			url, err := image.GetStorageURL()
			if err != nil {
				fmt.Println(err)
				return
			}

			director := func(req *http.Request) {
				req.URL = url
				req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
				req.Host = url.Host
			}

			proxy := &httputil.ReverseProxy{Director: director}
			proxy.ServeHTTP(w, r)
		} else {
			fmt.Println("404") //TODO
		}
	})

	http.Handle("/", r)
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
