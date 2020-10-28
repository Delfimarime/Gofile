package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type Publisher interface {
	publish(configuration Configuration, array []string) ([]bool, []int, [][]byte)
}

type PublisherImpl struct {
}

func (instance *PublisherImpl) publish(configuration Configuration, array []string) ([]bool, []int, [][]byte) {

	if len(array) == 0 {
		panic("no request to be sent")
	}

	if len(array) == 1 {
		v1, v2, v3 := instance.send(newClient(configuration), configuration, array)
		sent, statusCode, content := v1[0], v2[0], v3[0]
		return []bool{sent}, []int{statusCode}, [][]byte{content}
	}

	return instance.send(newClient(configuration), configuration, array)
}

func (instance *PublisherImpl) send(client http.Client, configuration Configuration, filename []string) ([]bool, []int, [][]byte) {

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

	if configuration.Verbose {
		fmt.Println("<----------------- Send " + filename + " through binary body ----------------->\n")
	}

	data, err := os.Open(filename)

	if configuration.Verbose {
		fmt.Println("Open:" + filename)
	}

	if err != nil {
		if configuration.Verbose {
			fmt.Println("Cannot Open:" + filename + " , cause:")
			fmt.Println(err)
		}
		return false, -1, nil
	}

	req, err := http.NewRequest("POST", configuration.Endpoint, data)

	if configuration.Verbose {
		fmt.Println("Building Http Request for:" + filename)
	}

	if err != nil {

		if configuration.Verbose {
			fmt.Println("Cannot build Http Request for:" + filename + ", cause:")
			fmt.Println(err)
		}

		return false, -1, nil
	}

	if len(configuration.Username) > 0 {

		req.SetBasicAuth(configuration.Username, configuration.Password)

		if configuration.Verbose {
			fmt.Println("Username & Password detected for basic authentication on " + configuration.Endpoint + " for file:" + filename)
		}

	}

	resp, err := client.Do(req)

	if err != nil {

		if configuration.Verbose {
			fmt.Println("Failed Http Request submission on " + configuration.Endpoint + " for:" + filename + ", cause:")
			fmt.Println(err)
		}

		return false, -1, nil
	}

	if configuration.Verbose {
		fmt.Println("Http Request submitted for:" + filename)
	}

	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {

		if configuration.Verbose {
			fmt.Println("Cannot ready Response for:" + filename + ", cause:")
			fmt.Println(err)
		}

		fmt.Println("<----------------- Send " + filename + " through binary body ----------------->\n")

		return resp.StatusCode == 200 || resp.StatusCode == 201, resp.StatusCode, nil
	}

	fmt.Println("<----------------- Send " + filename + " through binary body ----------------->\n")

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

func newClient(configuration Configuration) http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	timeout := configuration.Timeout

	if timeout < 0 {
		timeout = 10
	}

	return http.Client{Transport: transport, Timeout: time.Duration(timeout) * time.Second}
}
