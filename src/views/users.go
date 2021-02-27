package views

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/zer0tonin/senketsu/src/model"
)

func UsersHandler(r *mux.Router) {
	templates := make(map[string]*template.Template)
	templates["user"] = template.Must(template.ParseFiles("./templates/base.html", "./templates/user.html"))
	templates["users"] = template.Must(template.ParseFiles("./templates/base.html", "./templates/users.html"))
	templates["auth"] = template.Must(template.ParseFiles("./templates/base.html", "./templates/auth.html"))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		users, err := model.S.UserRepository.List(r.Context())
		if err != nil {
			fmt.Println(err)
			serveError(w, r, 500, "Failed to fetch users")
			return
		}
		err = templates["users"].Execute(
			w,
			map[string]interface{}{
				"users": users,
			},
		)
		if err != nil {
			fmt.Println(err)
			serveError(w, r, 500, "Failed to render")
		}
	})

	r.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		err := templates["auth"].Execute(
			w,
			map[string]interface{}{
			},
		)
		if err != nil {
			fmt.Println(err)
			serveError(w, r, 500, "Failed to render")
		}
	}).Methods("GET")

	r.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		// TODO: auth and set cookie
	}).Methods("POST")

	r.HandleFunc("/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user, err := model.S.UserRepository.Get(r.Context(), vars["id"])
		if err != nil {
			fmt.Println(err)
			serveError(w, r, 500, "Failed to fetch user")
			return
		}
		if user == nil {
			fmt.Println("Image not found")
			serveError(w, r, 404, "Image not found")
			return
		}
		err = templates["user"].Execute(
			w,
			map[string]interface{}{
				"user": user,
			},
		)
		if err != nil {
			fmt.Println(err)
			serveError(w, r, 500, "Failed to render")
		}
	})
}
