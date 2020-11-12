package views

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httputil"
	"regexp"

	"github.com/gorilla/mux"

	"github.com/zer0tonin/senketsu/src/model"
)

func BaseHandler() *mux.Router {
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
		w.Header().Set("Content-Type", "text/css")
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

	return r
}
