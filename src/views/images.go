package views

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/zer0tonin/senketsu/src/model"
)

func ImagesHandler(r *mux.Router) {
	templates := make(map[string]*template.Template)
	templates["image"] = template.Must(template.ParseFiles("./templates/base.html", "./templates/image.html"))
	templates["images"] = template.Must(template.ParseFiles("./templates/base.html", "./templates/images.html"))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		images, err := model.S.ImageRepository.List(r.Context())
		if err != nil {
			fmt.Println(err)
			serveError(w, r, 500, "Failed to fetch images")
			return
		}
		err = templates["images"].Execute(
			w,
			map[string]interface{}{
				"images": images,
			},
		)
		if err != nil {
			fmt.Println(err)
			serveError(w, r, 500, "Failed to render")
		}
	})

	r.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		image, err := model.S.ImageRepository.Get(r.Context(), vars["id"])
		if err != nil {
			fmt.Println(err)
			serveError(w, r, 500, "Failed to fetch image")
			return
		}
		if image == nil {
			fmt.Println("Image not found")
			serveError(w, r, 404, "Image not found")
			return
		}
		err = templates["image"].Execute(
			w,
			map[string]interface{}{
				"image": image,
			},
		)
		if err != nil {
			fmt.Println(err)
			serveError(w, r, 500, "Failed to render")
		}
	})
}
