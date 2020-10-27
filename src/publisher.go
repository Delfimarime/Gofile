package main

import (
	"fmt"
	"net/http"
	"time"
)

type Publisher interface {
	send(req http.Request) bool
	sendMany(array []http.Request) []bool
}

type PublisherImpl struct {
}

func (instance *PublisherImpl) send(req http.Request) bool {
	client := http.Client{}
	return instance.submit(client, &req)
}

func (instance *PublisherImpl) sendMany(array []http.Request) []bool {

	if len(array) == 0 {
		panic("no response to be sent")
	}

	if len(array) == 1 {
		return []bool{instance.send(array[0])}
	}

	answers := make([]bool, 0)

	client := http.Client{Timeout: time.Second * 20}

	for _, req := range array {
		answers = append(answers, instance.submit(client, &req))
	}

	return answers
}

func (instance *PublisherImpl) submit(client http.Client, req *http.Request) bool {

	resp, err := client.Do(req)

	if err != nil {
		return false
	}

	fmt.Println(resp)

	return resp.Status == "200 OK"
}
