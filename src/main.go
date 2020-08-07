package main

import (
	"fmt"
	"html/template"
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

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
