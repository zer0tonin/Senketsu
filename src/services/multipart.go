package services

import (
	"mime"
	"mime/multipart"
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
	_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		errs = append(errs, err)
		return
	}

	mr := multipart.NewReader(r.Body, params["boundary"])
	form, err := mr.ReadForm(m.maxMem) // 128MiB
	if err != nil {
		errs = append(errs, err)
		return
	}

	for _, fileHeaders := range form.File {
		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()
			if err != nil {
				errs = append(errs, err)
			} else {
				image := &model.Image{
					Extension: "gif", // FIXME: need to validate file types
					Reader:    file,
					Size:      fileHeader.Size,
				}
				result = append(result, image)
			}
		}
	}

	return
}
