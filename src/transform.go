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

	req, _, writer := toResponse(&content[0], configuration, index)

	if index > -1 {
		array := content[1:]

		index = 1

		for _, each := range array {
			attach(&writer, configuration, index, &each)
			index += 1
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

func attach(writer *multipart.Writer, configuration Configuration, index int, file *os.File) {

	var err error
	var fw io.Writer

	property := ""

	if len(configuration.Attribute) > 0 {
		property = configuration.Attribute
	}

	if len(property) == 0 {
		property = "file"
	}

	if index >= 0 {
		property = property + "[]"
	}

	if fw, err = writer.CreateFormFile(property, file.Name()); err != nil {
		panic(err)
	}

	if _, err = io.Copy(fw, file); err != nil {
		panic(err)
	}

}

func toResponse(file *os.File, configuration Configuration, index int) (http.Request, bytes.Buffer, multipart.Writer) {

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	attach(writer, configuration, index, file)

	req, err := http.NewRequest("POST", configuration.Endpoint, &buffer)

	if err != nil {
		panic(err)
	}
	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if len(configuration.Username) > 0 {
		req.SetBasicAuth(configuration.Username, configuration.Password)
	}

	return *req, buffer, *writer
}
