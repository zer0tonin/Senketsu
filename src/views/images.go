package views

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/zer0tonin/senketsu/src/model"
)

func ImagesHandler() *mux.Router {
	r := mux.NewRouter()
	templates := make(map[string]*template.Template)
	templates["image"] = template.Must(template.ParseFiles("./templates/base.html", "./templates/image.html"))
	templates["images"] = template.Must(template.ParseFiles("./templates/base.html", "./templates/images.html"))

	r.HandleFunc("/images/", func(w http.ResponseWriter, r *http.Request) {
		images, err := model.S.ImageRepository.List(r.Context())
		if err != nil {
			fmt.Println(err)
			return
		}
		templates["images"].Execute(
			w,
			map[string]interface{}{
				"images": images,
			},
		)
	})

	r.HandleFunc("/images/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		image, err := model.S.ImageRepository.Get(r.Context(), vars["id"])
		if err != nil {
			fmt.Println(err)
			return
		}
		templates["image"].Execute(
			w,
			map[string]interface{}{
				"image": image,
			},
		)
	})

	return r
}
