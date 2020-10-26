package main

import (
	"net/http"
)

type Sender interface {
	send(req http.Request) bool
	sendMany(array []http.Request) []bool
}

type DefaultSender struct {
}

func (instance *DefaultSender) send(req http.Request) bool {
	client := http.Client{}
	return instance.submit(client, &req)
}

func (instance *DefaultSender) sendMany(array []http.Request) []bool {

	if len(array) == 0 {
		panic("no response to be sent")
	}

	if len(array) == 1 {
		return []bool{instance.send(array[0])}
	}

	answers := make([]bool, 0)

	client := http.Client{}

	for _, req := range array {
		answers = append(answers, instance.submit(client, &req))
	}

	return answers
}

func (instance *DefaultSender) submit(client http.Client, req *http.Request) bool {
	resp, err := client.Do(req)

	if err != nil {
		return false
	}

	return resp.Status == "200 OK"
}
