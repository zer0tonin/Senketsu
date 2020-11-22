package services

import (
	"net/http"

	"github.com/zer0tonin/senketsu/src/model"
)

type MultipartParser struct {
	maxMem int64
}

func NewMultipartParser(maxMem int64) *MultipartParser {
	return &MultipartParser{
		maxMem: maxMem,
	}
}

func (m *MultipartParser) ParseForm(r *http.Request) (result []*model.Image, errs []error) {
	err := r.ParseMultipartForm(m.maxMem)
	if err != nil {
		errs = append(errs, err)
		return
	}
	form := r.MultipartForm
	data := r.PostForm

	for _, fileHeaders := range form.File { // FIXME: files need to be ordered in some way to match tags or idk
		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()
			if err != nil {
				errs = append(errs, err)
			} else {
				image := &model.Image{
					Extension: "gif", // FIXME: need to validate file types
					Reader:    file,
					Size:      fileHeader.Size,
					Tags:      data["tags"],
				}
				result = append(result, image)
			}
		}
	}

	return
}
