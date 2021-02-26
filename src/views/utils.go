package views

import (
	"fmt"
	"html/template"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Printf("%s %s\n", r.Method, r.RequestURI)
        next.ServeHTTP(w, r)
    })
}

var errorTemplate = template.Must(
	template.ParseFiles("./templates/base.html", "./templates/error.html"),
)

func serveError(w http.ResponseWriter, r *http.Request, code int, message string) {
	fmt.Printf("%d on %s %s\n", code, r.Method, r.RequestURI)
	w.WriteHeader(code)
	errorTemplate.Execute(
		w,
		map[string]interface{}{
			"message": message,
		},
	)
}
