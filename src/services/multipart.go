package services

import (
	"mime"
	"mime/multipart"
	"net/http"
)

type MultipartParser struct {
	maxMem int64
}

func NewMultipartParser(maxMem int64) *MultipartParser {
	return &MultipartParser{
		maxMem: maxMem,
	}
}

func (m *MultipartParser) ParseForm(r *http.Request) (result []*multipart.FileHeader, err error) {
	_, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil {
		return
	}

	mr := multipart.NewReader(r.Body, params["boundary"])
	form, err := mr.ReadForm(m.maxMem) // 128MiB
	if err != nil {
		return
	}

	for _, fileHeaders := range form.File {
		result = append(result, fileHeaders...)
	}
	return
}
