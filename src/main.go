package main

import (
	"fmt"
	"html/template"
	"mime"
	"mime/multipart"
	"net/http"
)

type HelloData struct {
	Name string
}

func main() {
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := HelloData{
			Name: "Alice",
		}
		tmpl.Execute(w, data)
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
