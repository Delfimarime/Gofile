package main

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type Transformer interface {
	transform(content []os.File, configuration Configuration) []http.Request
}

type DefaultTransformer struct {
}

func (instance *DefaultTransformer) transform(content []os.File, configuration Configuration) []http.Request {
	if configuration.Compact {
		return toMono(content, configuration)
	} else {
		return toFlux(content, configuration)
	}
}

func toMono(content []os.File, configuration Configuration) []http.Request {

	if len(content) == 0 {
		panic(errors.New("at least one (1) file is expected"))
	}

	var index = -1

	if len(content) > 0 {
		index = 0
	}

	req, _, writer := toResponse(&content[0], configuration)

	if index > -1 {
		array := content[1:]

		for _, each := range array {
			setFile(&writer, configuration, &each)
		}
	}

	err := writer.Close()

	if err != nil {
		panic(err)
	}

	return []http.Request{req}
}

func toFlux(content []os.File, configuration Configuration) []http.Request {

	if len(content) < 2 {
		return toMono(content, configuration)
	}

	array := make([]http.Request, 0)
	for _, each := range content {
		array = append(array, toMono([]os.File{each}, configuration)...)
	}
	return array
}

func setFile(writer *multipart.Writer, configuration Configuration, file *os.File) {

	var err error
	var fw io.Writer

	property := ""

	if len(configuration.Attribute) > 0 {
		property = configuration.Attribute
	}

	if len(property) == 0 {
		property = "file"
	}

	if fw, err = writer.CreateFormFile(property, file.Name()); err != nil {
		panic(err)
	}

	if _, err = io.Copy(fw, file); err != nil {
		panic(err)
	}

}

func toResponse(file *os.File, configuration Configuration) (http.Request, bytes.Buffer, multipart.Writer) {

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	setFile(writer, configuration, file)

	req, err := http.NewRequest("POST", configuration.Endpoint, &body)

	if err != nil {
		panic(err)
	}

	req.Header.Del("Accept")
	req.Header.Del("Content-Type")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if len(configuration.Username) > 0 {
		req.SetBasicAuth(configuration.Username, configuration.Password)
	}

	return *req, body, *writer
}
