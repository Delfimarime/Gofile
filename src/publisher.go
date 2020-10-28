package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type Publisher interface {
	send(configuration Configuration, file string) (bool, int, []byte)
	sendMany(configuration Configuration, array []string) ([]bool, []int, [][]byte)
}

type PublisherImpl struct {
}

func (instance *PublisherImpl) send(configuration Configuration, filename string) (bool, int, []byte) {
	client := http.Client{}
	return instance.submitSingle(client, configuration, filename)
}

func (instance *PublisherImpl) sendMany(configuration Configuration, array []string) ([]bool, []int, [][]byte) {

	if len(array) == 0 {
		panic("no request to be sent")
	}

	if len(array) == 1 {
		sent, statusCode, content := instance.send(configuration, array[0])
		return []bool{sent}, []int{statusCode}, [][]byte{content}
	}

	client := http.Client{}

	return instance.submitMany(client, configuration, array)
}

func (instance *PublisherImpl) submitSingle(client http.Client, configuration Configuration, filename string) (bool, int, []byte) {

	withAttribute := len(configuration.Attribute) == 0

	if !withAttribute && configuration.Compact {
		panic(errors.New("compact isn't supported"))
	}

	if !withAttribute {
		return instance.sendBinary(client, configuration, filename)
	}

	return instance.sendForm(client, configuration, []string{filename})
}

func (instance *PublisherImpl) submitMany(client http.Client, configuration Configuration, filename []string) ([]bool, []int, [][]byte) {

	withAttribute := len(configuration.Attribute) > 0

	if !withAttribute && configuration.Compact {
		panic(errors.New("compact isn't supported"))
	}

	if !withAttribute {

		status := make([]bool, 0)
		array := make([][]byte, 0)
		statusCode := make([]int, 0)

		for _, each := range filename {
			v1, v2, v3 := instance.sendBinary(client, configuration, each)
			array = append(array, v3)
			status = append(status, v1)
			statusCode = append(statusCode, v2)
		}

		return status, statusCode, array

	}

	if configuration.Compact {
		v1, v2, v3 := instance.sendForm(client, configuration, filename)
		return []bool{v1}, []int{v2}, [][]byte{v3}
	}

	status := make([]bool, 0)
	array := make([][]byte, 0)
	statusCode := make([]int, 0)

	for _, each := range filename {
		v1, v2, v3 := instance.sendForm(client, configuration, []string{each})
		array = append(array, v3)
		status = append(status, v1)
		statusCode = append(statusCode, v2)
	}

	return status, statusCode, array
}

func (instance *PublisherImpl) sendForm(client http.Client, configuration Configuration, filename []string) (bool, int, []byte) {

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	for _, each := range filename {
		setFile(writer, configuration.Attribute, each)
	}

	err := writer.Close()

	if err != nil {
		panic(err)
	}

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

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	return resp.StatusCode == 200 || resp.StatusCode == 201, resp.StatusCode, content
}

func (instance *PublisherImpl) sendBinary(client http.Client, configuration Configuration, filename string) (bool, int, []byte) {

	data, err := os.Open(filename)

	if err != nil {
		panic(err)
		//return false , nil, make([]byte,0)
	}

	req, err := http.NewRequest("POST", configuration.Endpoint, data)

	if err != nil {
		panic(err)
		//return false , nil, make([]byte,0)
	}

	if len(configuration.Username) > 0 {
		req.SetBasicAuth(configuration.Username, configuration.Password)
	}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
		//return false , nil, make([]byte,0)
	}

	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return resp.StatusCode == 200 || resp.StatusCode == 201, resp.StatusCode, nil
	}

	return resp.StatusCode == 200 || resp.StatusCode == 201, resp.StatusCode, content

}

func setFile(writer *multipart.Writer, attr string, filename string) {

	var err error
	var fw io.Writer

	file, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	if fw, err = writer.CreateFormFile(attr, file.Name()); err != nil {
		panic(err)
	}

	if _, err = io.Copy(fw, file); err != nil {
		panic(err)
	}

}
