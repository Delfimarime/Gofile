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
	transformer     Transformer
	discoveryClient DiscoveryClient
}

func (instance *GoEngine) Run(configuration Configuration) {

	if instance.sender == nil {
		panic(errors.New("sender mustn't be null"))
	}

	if instance.transformer == nil {
		panic(errors.New("transformer mustn't be null"))
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

	req := instance.transformer.transform(files, configuration)

	if len(req) == 0 {
		panic("no request to be submitted:" + configuration.File)
	}

	resp := instance.sender.sendMany(req)

	fmt.Println("----------------------------- REPORT -----------------------------")

	if configuration.Compact {
		fmt.Println(asMono(files, resp))
	} else {
		fmt.Println(asFlux(files, resp))
	}

	fmt.Println("----------------------------- REPORT -----------------------------")

	var isValid = true

	for index := range files {
		isValid = resp[index]

		if (configuration.Strategy == AtLeastOne && isValid) || configuration.Strategy != EveryFile {
			return
		} else if configuration.Strategy != AtLeastOne && !isValid {
			fmt.Println("\nexpected every file to be upload but an error has occurred")
			os.Exit(1)
		}

	}

}

func (instance *GoEngine) SetSender(sender Publisher) {
	if sender != nil {
		instance.sender = sender
	}
}

func asFlux(files []os.File, resp []bool) string {
	s := ""

	for index, each := range files {
		s += each.Name() + " => " + strconv.FormatBool(resp[index]) + "\n"
	}

	return s
}

func asMono(files []os.File, resp []bool) string {
	s := ""

	for index, each := range files {
		if index == 0 {
			s = each.Name()
		} else {
			s += " , " + each.Name()
		}
	}

	s += " => " + strconv.FormatBool(resp[0])

	return s
}

func (instance *GoEngine) SetTransformer(transformer Transformer) {
	if transformer != nil {
		instance.transformer = transformer
	}
}

func (instance *GoEngine) SetDiscoveryClient(discoveryClient DiscoveryClient) {
	if discoveryClient != nil {
		instance.discoveryClient = discoveryClient
	}
}
