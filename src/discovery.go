package main

import (
	"os"
	"path/filepath"
)

type DiscoveryClient interface {
	getFiles(path string) []os.File
}

type BasicDiscoveryClient struct {
}

func (instance *BasicDiscoveryClient) getFiles(path string) []os.File {

	if len(path) == 0 {
		return make([]os.File, 0)
	}

	file := openFile(path)

	fi, err := file.Stat()

	if err != nil {
		panic(err)
	}

	if !fi.IsDir() {
		return []os.File{*file}
	}

	var children []string

	err = filepath.Walk(path, addChildren(&children))

	if err != nil {
		panic(err)
	}

	array := make([]os.File, 0)

	for _, child := range children {
		array = append(array, *openFile(child))
	}

	return array
}
