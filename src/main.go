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
	fmt.Println("Hello world")
	tmpl := template.Must(template.ParseFiles("/app/templates/index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := HelloData{
			Name: "Alice",
		}
		tmpl.Execute(w, data)
	})
	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
