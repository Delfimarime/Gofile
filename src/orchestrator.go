package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	EveryFile  string = "EVERY_FILE"
	AtLeastOne string = "AT_LEAST_ONE"
)

type Configuration struct {
	Timeout   int
	Verbose   bool
	Compact   bool
	File      string
	Endpoint  string
	Username  string
	Password  string
	Attribute string
	Strategy  string
}

type GoEngine struct {
	sender          Publisher
	discoveryClient DiscoveryClient
}

func (instance *GoEngine) Run(configuration Configuration) {

	if instance.sender == nil {
		panic(errors.New("sender mustn't be null"))
	}

	if instance.discoveryClient == nil {
		panic(errors.New("discovery client mustn't be null"))
	}

	if len(configuration.File) == 0 {
		panic(errors.New("cannot delete directory/File:" + configuration.File))
		return
	}

	if len(configuration.Endpoint) == 0 {
		panic(errors.New("endpoint mustn't be empty"))
		return
	}

	if !(strings.HasPrefix(configuration.Endpoint, "http://") || strings.HasPrefix(configuration.Endpoint, "https://")) {
		panic(errors.New("endpoint start with http:// or https://"))
		return
	}

	var hasUsername = len(configuration.Username) > 0
	var hasPassword = len(configuration.Password) > 0

	if hasUsername != hasPassword {
		if hasUsername && !hasPassword {
			panic(errors.New("password is required"))
		} else {
			panic(errors.New("username is required"))
		}
		return
	}

	files := instance.discoveryClient.getFiles(configuration.File)

	if len(files) == 0 {
		panic("no file found for:" + configuration.File)
	}

	sent, statusCode, content := instance.sender.sendMany(configuration, files)

	everyFileSent := analyze(configuration, files, sent, statusCode, content)

	if (configuration.Strategy == AtLeastOne && everyFileSent) || configuration.Strategy != EveryFile {
		return
	} else if configuration.Strategy != AtLeastOne && !everyFileSent {
		fmt.Println("\nexpected every file to be upload but an error has occurred")
		os.Exit(1)
	}

}

func analyze(configuration Configuration, files []string, sent []bool, statusCode []int, content [][]byte) bool {
	fmt.Println("\n----------------------------- REPORT -----------------------------")
	everyFileSent := true

	if configuration.Compact {

		for _, each := range files {
			fmt.Println("filename   :	" + each)
		}

		fmt.Println("uploaded   :	" + strconv.FormatBool(sent[0]))
		fmt.Println("status code:	" + strconv.Itoa(statusCode[0]))

		if configuration.Verbose {
			fmt.Println("content    :	" + string(content[0]))
		}

		everyFileSent = everyFileSent && sent[0]

	} else {

		for index := range files {

			fmt.Println("filename   :	" + files[index])
			fmt.Println("uploaded   :	" + strconv.FormatBool(sent[index]))
			fmt.Println("status code:	" + strconv.Itoa(statusCode[index]))

			if configuration.Verbose {
				fmt.Println("content    :	" + string(content[index]))
			}

			fmt.Print("\n\n")

			if everyFileSent && !sent[index] {
				everyFileSent = false
			}

		}

	}
	fmt.Println("\n----------------------------- REPORT -----------------------------")
	return everyFileSent
}

func (instance *GoEngine) SetSender(sender Publisher) {
	if sender != nil {
		instance.sender = sender
	}
}

func (instance *GoEngine) SetDiscoveryClient(discoveryClient DiscoveryClient) {
	if discoveryClient != nil {
		instance.discoveryClient = discoveryClient
	}
}
