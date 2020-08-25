package main

import (
	"fmt"
	"html/template"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/spf13/viper"
)

func config() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config"))
	} else {
		fmt.Println(viper.GetString("test"))
	}
}

func serve() {
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/index.css")
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			fmt.Println(err)
			return
		}
		mr := multipart.NewReader(r.Body, params["boundary"])
		part, err := mr.NextPart()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(part.FileName())
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func main() {
	config()
	serve()
}
