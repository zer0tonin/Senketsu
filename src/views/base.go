package views

import (
	"fmt"
	"html/template"
	"os"
	"net/http"
	"net/http/httputil"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"

	"github.com/zer0tonin/senketsu/src/model"
)

func BaseHandler() http.Handler {
	templates := make(map[string]*template.Template)
	templates["index"] = template.Must(template.ParseFiles("./templates/base.html", "./templates/index.html"))
	templates["upload"] = template.Must(template.ParseFiles("./templates/base.html", "./templates/upload.html"))

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tags, err := model.S.TagRepository.List(r.Context())
		if err != nil {
			fmt.Println(err)
			serveError(w, r, 500, "Failed to fetch tag list")
			return
		}
		err = templates["index"].Execute(
			w,
			map[string]interface{}{
				"tags": tags,
			},
		)
		if err != nil {
			fmt.Println(err)
			serveError(w, r, 500, "Failed to render")
		}
	})

	r.HandleFunc("/index.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		http.ServeFile(w, r, "./templates/index.css")
	})

	r.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/javascript")
		http.ServeFile(w, r, "./templates/index.js")
	})

	r.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		images, errs := model.NewImageFromRequest(r.Context(), r)
		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Println(err)
			}
			serveError(w, r, 400, "Failed to upload image")
			return
		}

		err := templates["upload"].Execute(
			w,
			map[string]interface{}{
				"images": images,
				"errs":   errs,
			},
		)
		if err != nil {
			fmt.Println(err)
			serveError(w, r, 500, "Failed to render")
		}
	})

	imagesRegexp := regexp.MustCompile(`.+\/(.+)\.gif`)
	r.HandleFunc("/files/{path}", func(w http.ResponseWriter, r *http.Request) {
		/*
		Reverse proxy route for S3
		*/
		matches := imagesRegexp.FindStringSubmatch(r.URL.Path)
		if len(matches) == 2 {
			image, err := model.S.ImageRepository.Get(r.Context(), matches[1])
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(500)
				return
			}
			if image == nil {
				fmt.Println("Image not found")
				w.WriteHeader(404)
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
			fmt.Println("Image not found")
			w.WriteHeader(404)
		}
	})

	UsersHandler(r.PathPrefix("/users").Subrouter())
	ImagesHandler(r.PathPrefix("/images").Subrouter())

	return handlers.LoggingHandler(
		os.Stdout,
		handlers.RecoveryHandler()(r),
	)
}
